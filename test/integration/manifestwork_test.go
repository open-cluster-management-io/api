package integration

import (
	"context"
	"fmt"

	workv1 "open-cluster-management.io/api/work/v1"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

var _ = ginkgo.Describe("ManifestWork API test", func() {
	var manifestWorkName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		manifestWorkName = fmt.Sprintf("manifestwork-%s", suffix)
	})

	ginkgo.It("Create a ManifestWork with empty executor identity", func() {
		work := &workv1.ManifestWork{
			ObjectMeta: metav1.ObjectMeta{
				Name: manifestWorkName,
			},
		}
		_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
			Create(context.TODO(), work, metav1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	invalidSubjectType := "my-invalid-service-account-name"
	ginkgo.It("Create a ManifestWork with an non-empty identity", func() {
		work := &workv1.ManifestWork{
			ObjectMeta: metav1.ObjectMeta{
				Name: manifestWorkName,
			},
			Spec: workv1.ManifestWorkSpec{
				Executor: &workv1.ManifestWorkExecutor{
					Subject: workv1.ManifestWorkExecutorSubject{
						Type: workv1.ManifestWorkExecutorSubjectType(invalidSubjectType),
					},
				},
			},
		}
		_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
			Create(context.TODO(), work, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("Create a ManifestWork with an non-empty identity", func() {
		work := &workv1.ManifestWork{
			ObjectMeta: metav1.ObjectMeta{
				Name: manifestWorkName,
			},
			Spec: workv1.ManifestWorkSpec{
				Executor: &workv1.ManifestWorkExecutor{
					Subject: workv1.ManifestWorkExecutorSubject{
						Type: workv1.ExecutorSubjectTypeServiceAccount,
						ServiceAccount: &workv1.ManifestWorkSubjectServiceAccount{
							Namespace: "default",
							Name:      "----",
						},
					},
				},
			},
		}
		_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
			Create(context.TODO(), work, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

})
