package api

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

	ginkgo.Context("ManifestWork executor", func() {
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

	ginkgo.Context("ManifestWork update strategy", func() {
		ginkgo.It("Create a ManifestWork with no configs", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					ManifestConfigs: []workv1.ManifestConfigOption{},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("Create a ManifestWork with empty update strategy", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					ManifestConfigs: []workv1.ManifestConfigOption{
						{
							ResourceIdentifier: workv1.ResourceIdentifier{
								Resource:  "foo",
								Name:      "test",
								Namespace: "testns",
							},
							UpdateStrategy: &workv1.UpdateStrategy{},
						},
					},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			work, err = hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Get(context.TODO(), work.Name, metav1.GetOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(work.Spec.ManifestConfigs[0].UpdateStrategy.Type).Should(gomega.Equal(workv1.UpdateStrategyTypeUpdate))
		})

		ginkgo.It("Create a ManifestWork with ssa update strategy and default manager", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					ManifestConfigs: []workv1.ManifestConfigOption{
						{
							ResourceIdentifier: workv1.ResourceIdentifier{
								Resource:  "foo",
								Name:      "test",
								Namespace: "testns",
							},
							UpdateStrategy: &workv1.UpdateStrategy{
								Type:            workv1.UpdateStrategyTypeServerSideApply,
								ServerSideApply: &workv1.ServerSideApplyConfig{},
							},
						},
					},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			work, err = hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Get(context.TODO(), work.Name, metav1.GetOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(work.Spec.ManifestConfigs[0].UpdateStrategy.ServerSideApply.FieldManager).Should(gomega.Equal("work-agent"))
		})

		ginkgo.It("Create a ManifestWork with ssa update strategy and invalid manager", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					ManifestConfigs: []workv1.ManifestConfigOption{
						{
							ResourceIdentifier: workv1.ResourceIdentifier{
								Resource:  "foo",
								Name:      "test",
								Namespace: "testns",
							},
							UpdateStrategy: &workv1.UpdateStrategy{
								Type: workv1.UpdateStrategyTypeServerSideApply,
								ServerSideApply: &workv1.ServerSideApplyConfig{
									FieldManager: "another",
								},
							},
						},
					},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).To(gomega.HaveOccurred())

			work.Spec.ManifestConfigs[0].UpdateStrategy.ServerSideApply.FieldManager = "work-agent-another"

			_, err = hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("ManifestWork feedback rules", func() {
		ginkgo.It("feedback rule type must be set", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					ManifestConfigs: []workv1.ManifestConfigOption{
						{
							ResourceIdentifier: workv1.ResourceIdentifier{
								Resource:  "foo",
								Name:      "test",
								Namespace: "testns",
							},
							FeedbackRules: []workv1.FeedbackRule{
								{},
							},
						},
					},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("set feedback rule", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					ManifestConfigs: []workv1.ManifestConfigOption{
						{
							ResourceIdentifier: workv1.ResourceIdentifier{
								Resource:  "foo",
								Name:      "test",
								Namespace: "testns",
							},
							FeedbackRules: []workv1.FeedbackRule{
								{Type: workv1.WellKnownStatusType},
							},
						},
					},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})
	})
})
