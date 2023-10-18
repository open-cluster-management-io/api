package api

import (
	"context"
	"path/filepath"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	addonv1alpha1client "open-cluster-management.io/api/client/addon/clientset/versioned"
	clusterv1client "open-cluster-management.io/api/client/cluster/clientset/versioned"
	operatorclientset "open-cluster-management.io/api/client/operator/clientset/versioned"
	workclientset "open-cluster-management.io/api/client/work/clientset/versioned"
)

var testEnv *envtest.Environment
var testNamespace string
var kubernetesClient kubernetes.Interface
var hubWorkClient workclientset.Interface
var hubClusterClient clusterv1client.Interface
var hubAddonClient addonv1alpha1client.Interface
var operatorClient operatorclientset.Interface

func TestIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "API Validation Integration Suite")
}

var _ = ginkgo.BeforeSuite(func(done ginkgo.Done) {
	ginkgo.By("bootstrapping test environment")

	// start a kube-apiserver
	testEnv = &envtest.Environment{
		ErrorIfCRDPathMissing: true,
		CRDDirectoryPaths: []string{
			filepath.Join(".", "work", "v1", "0000_00_work.open-cluster-management.io_manifestworks.crd.yaml"),
			filepath.Join(".", "work", "v1", "0000_01_work.open-cluster-management.io_appliedmanifestworks.crd.yaml"),
			filepath.Join(".", "work", "v1alpha1"),
			filepath.Join(".", "cluster", "v1"),
			filepath.Join(".", "cluster", "v1beta1", "0000_02_clusters.open-cluster-management.io_placements.crd.yaml"),
			filepath.Join(".", "cluster", "v1beta1", "0000_03_clusters.open-cluster-management.io_placementdecisions.crd.yaml"),
			filepath.Join(".", "cluster", "v1beta2"),

			filepath.Join(".", "cluster", "v1alpha1", "0000_02_clusters.open-cluster-management.io_clusterclaims.crd.yaml"),
			filepath.Join(".", "cluster", "v1alpha1", "0000_05_clusters.open-cluster-management.io_addonplacementscores.crd.yaml"),
			filepath.Join(".", "addon", "v1alpha1"),
			filepath.Join(".", "operator", "v1", "0000_00_operator.open-cluster-management.io_klusterlets.crd.yaml"),
			filepath.Join(".", "operator", "v1", "0000_01_operator.open-cluster-management.io_clustermanagers.crd.yaml"),
		},
	}

	cfg, err := testEnv.Start()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(cfg).ToNot(gomega.BeNil())

	hubWorkClient, err = workclientset.NewForConfig(cfg)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	hubClusterClient, err = clusterv1client.NewForConfig(cfg)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	hubAddonClient, err = addonv1alpha1client.NewForConfig(cfg)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	operatorClient, err = operatorclientset.NewForConfig(cfg)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	kubernetesClient, err = kubernetes.NewForConfig(cfg)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	testNamespace = "open-cluster-management-api-test-" + rand.String(5)
	_, err = kubernetesClient.CoreV1().Namespaces().
		Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: testNamespace,
			},
		}, metav1.CreateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	close(done)
}, 300)

var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("tearing down the test environment")

	// Skip if client wasn't instantiated
	if kubernetesClient != nil {
		err := kubernetesClient.CoreV1().Namespaces().
			Delete(context.TODO(), testNamespace, metav1.DeleteOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		err = testEnv.Stop()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}
})
