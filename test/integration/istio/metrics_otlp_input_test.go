//go:build istio

package istio

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	testutils "github.com/kyma-project/telemetry-manager/internal/utils/test"
	"github.com/kyma-project/telemetry-manager/test/testkit/assert"
	kitk8s "github.com/kyma-project/telemetry-manager/test/testkit/k8s"
	kitkyma "github.com/kyma-project/telemetry-manager/test/testkit/kyma"
	kitbackend "github.com/kyma-project/telemetry-manager/test/testkit/mocks/backend"
	"github.com/kyma-project/telemetry-manager/test/testkit/mocks/telemetrygen"
	"github.com/kyma-project/telemetry-manager/test/testkit/periodic"
	"github.com/kyma-project/telemetry-manager/test/testkit/suite"
)

var _ = Describe(suite.ID(), Label(suite.LabelGardener, suite.LabelIstio), Ordered, func() {
	const (
		metricProducer1Name = "metric-producer-1"
		metricProducer2Name = "metric-producer-2"
	)

	var (
		backendNs          = suite.ID()
		istiofiedBackendNs = suite.IDWithSuffix("istiofied")

		pipeline1Name    = suite.IDWithSuffix("1")
		pipeline2Name    = suite.IDWithSuffix("2")
		backend          *kitbackend.Backend
		istiofiedBackend *kitbackend.Backend
	)

	makeResources := func() []client.Object {
		var objs []client.Object

		objs = append(objs, kitk8s.NewNamespace(backendNs).K8sObject())
		objs = append(objs, kitk8s.NewNamespace(istiofiedBackendNs, kitk8s.WithIstioInjection()).K8sObject())

		// Mocks namespace objects
		backend = kitbackend.New(backendNs, kitbackend.SignalTypeMetrics)
		objs = append(objs, backend.K8sObjects()...)

		istiofiedBackend = kitbackend.New(istiofiedBackendNs, kitbackend.SignalTypeMetrics)
		objs = append(objs, istiofiedBackend.K8sObjects()...)

		metricPipeline := testutils.NewMetricPipelineBuilder().
			WithName(pipeline1Name).
			WithOTLPOutput(testutils.OTLPEndpoint(backend.Endpoint())).
			Build()
		objs = append(objs, &metricPipeline)

		metricPipelineIstiofiedBackend := testutils.NewMetricPipelineBuilder().
			WithName(pipeline2Name).
			WithOTLPOutput(testutils.OTLPEndpoint(istiofiedBackend.Endpoint())).
			Build()

		objs = append(objs, &metricPipelineIstiofiedBackend)

		// set peerauthentication to strict explicitly
		peerAuth := kitk8s.NewPeerAuthentication(kitbackend.DefaultName, istiofiedBackendNs)
		objs = append(objs, peerAuth.K8sObject(kitk8s.WithLabel("app", kitbackend.DefaultName)))

		// Create 2 deployments (with and without side-car) which would push the metrics to the metrics gateway.
		podSpec := telemetrygen.PodSpec(telemetrygen.SignalTypeMetrics)
		objs = append(objs,
			kitk8s.NewDeployment(metricProducer1Name, backendNs).WithPodSpec(podSpec).K8sObject(),
			kitk8s.NewDeployment(metricProducer2Name, istiofiedBackendNs).WithPodSpec(podSpec).K8sObject(),
		)

		return objs
	}

	Context("Istiofied and non-istiofied in-cluster backends", Ordered, func() {
		BeforeAll(func() {
			k8sObjects := makeResources()

			DeferCleanup(func() {
				Expect(kitk8s.DeleteObjects(k8sObjects...)).Should(Succeed())
				for _, resource := range k8sObjects {
					Eventually(func(g Gomega) {
						key := types.NamespacedName{Name: resource.GetName(), Namespace: resource.GetNamespace()}
						err := suite.K8sClient.Get(suite.Ctx, key, resource)
						g.Expect(apierrors.IsNotFound(err)).To(BeTrueBecause("Resource %s still exists", key))
					}, periodic.EventuallyTimeout, periodic.DefaultInterval).Should(Succeed())
				}
			})

			Expect(kitk8s.CreateObjects(GinkgoT(), k8sObjects...)).Should(Succeed())
		})

		It("Should have a running metric gateway deployment", func() {
			assert.DeploymentReady(GinkgoT(), kitkyma.MetricGatewayName)
		})

		It("Should have reachable backends", func() {
			assert.BackendReachable(GinkgoT(), backend)
			assert.BackendReachable(GinkgoT(), istiofiedBackend)
		})

		It("Should push metrics successfully", func() {
			assert.MetricsFromNamespaceDelivered(GinkgoT(), backend, backendNs, telemetrygen.MetricNames)
			assert.MetricsFromNamespaceDelivered(GinkgoT(), backend, istiofiedBackendNs, telemetrygen.MetricNames)

			assert.MetricsFromNamespaceDelivered(GinkgoT(), istiofiedBackend, backendNs, telemetrygen.MetricNames)
			assert.MetricsFromNamespaceDelivered(GinkgoT(), istiofiedBackend, istiofiedBackendNs, telemetrygen.MetricNames)

		})
	})
})
