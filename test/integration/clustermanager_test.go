package integration

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	operatorv1 "open-cluster-management.io/api/operator/v1"
)

var _ = ginkgo.Describe("ClusterManager API test", func() {
	var clusterManagerName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		clusterManagerName = fmt.Sprintf("cm-%s", suffix)
	})

	ginkgo.It("Create a cluster manager with empty spec", func() {
		clusterset := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{},
		}
		clusterset, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		gomega.Expect(clusterset.Spec.DeployOption.Mode).Should(gomega.Equal(operatorv1.InstallModeDefault))
	})

	ginkgo.It("Create a cluster manager with wrong install mode", func() {
		clusterset := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				DeployOption: operatorv1.ClusterManagerDeployOption{
					Mode: "WrongMode",
				},
			},
		}
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("Create a cluster manager with detached mode", func() {
		clusterset := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				DeployOption: operatorv1.ClusterManagerDeployOption{
					Mode:     operatorv1.InstallModeDetached,
					Detached: &operatorv1.DetachedClusterManagerConfiguration{},
				},
			},
		}
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())

		clusterset.Spec.DeployOption.Detached = &operatorv1.DetachedClusterManagerConfiguration{
			RegistrationWebhookConfiguration: operatorv1.WebhookConfiguration{
				Address: "test:test",
			},
			WorkWebhookConfiguration: operatorv1.WebhookConfiguration{
				Address: "test:test",
			},
		}
		_, err = operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())

		clusterset.Spec.DeployOption.Detached = &operatorv1.DetachedClusterManagerConfiguration{
			RegistrationWebhookConfiguration: operatorv1.WebhookConfiguration{
				Address: "localhost",
			},
			WorkWebhookConfiguration: operatorv1.WebhookConfiguration{
				Address: "localhost",
			},
		}
		clusterset, err = operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		gomega.Expect(clusterset.Spec.DeployOption.Detached.RegistrationWebhookConfiguration.Port).Should(gomega.Equal(int32(443)))
		gomega.Expect(clusterset.Spec.DeployOption.Detached.WorkWebhookConfiguration.Port).Should(gomega.Equal(int32(443)))
	})
})
