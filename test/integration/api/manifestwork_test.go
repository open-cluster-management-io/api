// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"time"

	workv1 "open-cluster-management.io/api/work/v1"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

		ginkgo.It("set well known status feedback rule", func() {
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

		ginkgo.It("set jsonpath feedback rule", func() {
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
								{
									Type: workv1.JSONPathsType,
									JsonPaths: []workv1.JsonPath{
										{Name: "Replica", Path: ".spec.replicas"},
										{Name: "StatusReplica", Path: ".status.replicas"},
									},
								},
							},
						},
					},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).
				Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("set feedback rule with same jsonpath name", func() {
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
								{
									Type: workv1.JSONPathsType,
									JsonPaths: []workv1.JsonPath{
										{Name: "Replica", Path: ".spec.replicas"},
										{Name: "Replica", Path: ".status.replicas"},
									},
								},
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
})

var _ = ginkgo.Describe("ManifestWork v1 Enhanced API test", func() {
	var manifestWorkName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		manifestWorkName = fmt.Sprintf("manifestwork-enhanced-%s", suffix)
	})

	ginkgo.AfterEach(func() {
		err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Delete(context.TODO(), manifestWorkName, metav1.DeleteOptions{})
		if !errors.IsNotFound(err) {
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		}
	})

	ginkgo.Context("ManifestWork workload validation", func() {
		ginkgo.It("should accept valid kubernetes manifests", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					Workload: workv1.ManifestsTemplate{
						Manifests: []workv1.Manifest{
							{
								RawExtension: runtime.RawExtension{
									Raw: []byte(`{
										"apiVersion": "v1",
										"kind": "ConfigMap",
										"metadata": {
											"name": "test-cm",
											"namespace": "default"
										},
										"data": {
											"key": "value"
										}
									}`),
								},
							},
						},
					},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("should handle multiple manifests", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					Workload: workv1.ManifestsTemplate{
						Manifests: []workv1.Manifest{
							{
								RawExtension: runtime.RawExtension{
									Raw: []byte(`{
										"apiVersion": "v1",
										"kind": "ConfigMap",
										"metadata": {
											"name": "test-cm-1"
										}
									}`),
								},
							},
							{
								RawExtension: runtime.RawExtension{
									Raw: []byte(`{
										"apiVersion": "v1",
										"kind": "Secret",
										"metadata": {
											"name": "test-secret-1"
										},
										"type": "Opaque"
									}`),
								},
							},
						},
					},
				},
			}

			manifestWork, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(manifestWork.Spec.Workload.Manifests)).Should(gomega.Equal(2))
		})
	})

	ginkgo.Context("ManifestWork delete options", func() {
		ginkgo.It("should reject empty propagation policy", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					DeleteOption: &workv1.DeleteOption{},
				},
			}

			_, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("Unsupported value: \"\""))
		})

		ginkgo.It("should accept valid propagation policy", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					DeleteOption: &workv1.DeleteOption{
						PropagationPolicy: workv1.DeletePropagationPolicyTypeForeground,
					},
				},
			}

			manifestWork, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(manifestWork.Spec.DeleteOption.PropagationPolicy).Should(gomega.Equal(workv1.DeletePropagationPolicyTypeForeground))
		})

		ginkgo.It("should accept custom propagation policy", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					DeleteOption: &workv1.DeleteOption{
						PropagationPolicy: workv1.DeletePropagationPolicyTypeOrphan,
					},
				},
			}

			manifestWork, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(manifestWork.Spec.DeleteOption.PropagationPolicy).Should(gomega.Equal(workv1.DeletePropagationPolicyTypeOrphan))
		})

		ginkgo.It("should handle TTL configuration", func() {
			ttl := int64(3600)
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					DeleteOption: &workv1.DeleteOption{
						PropagationPolicy:       workv1.DeletePropagationPolicyTypeForeground,
						TTLSecondsAfterFinished: &ttl,
					},
				},
			}

			manifestWork, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(*manifestWork.Spec.DeleteOption.TTLSecondsAfterFinished).Should(gomega.Equal(int64(3600)))
		})
	})

	ginkgo.Context("ManifestWork complex configuration", func() {
		ginkgo.It("should handle complete configuration with multiple features", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					Workload: workv1.ManifestsTemplate{
						Manifests: []workv1.Manifest{
							{
								RawExtension: runtime.RawExtension{
									Raw: []byte(`{
										"apiVersion": "apps/v1",
										"kind": "Deployment",
										"metadata": {
											"name": "test-deployment"
										},
										"spec": {
											"replicas": 2
										}
									}`),
								},
							},
						},
					},
					ManifestConfigs: []workv1.ManifestConfigOption{
						{
							ResourceIdentifier: workv1.ResourceIdentifier{
								Resource:  "deployments",
								Name:      "test-deployment",
								Namespace: "default",
							},
							UpdateStrategy: &workv1.UpdateStrategy{
								Type: workv1.UpdateStrategyTypeServerSideApply,
								ServerSideApply: &workv1.ServerSideApplyConfig{
									FieldManager: "work-agent-test",
									Force:        true,
								},
							},
							FeedbackRules: []workv1.FeedbackRule{
								{
									Type: workv1.JSONPathsType,
									JsonPaths: []workv1.JsonPath{
										{Name: "Replicas", Path: ".spec.replicas"},
										{Name: "ReadyReplicas", Path: ".status.readyReplicas"},
									},
								},
							},
						},
					},
					Executor: &workv1.ManifestWorkExecutor{
						Subject: workv1.ManifestWorkExecutorSubject{
							Type: workv1.ExecutorSubjectTypeServiceAccount,
							ServiceAccount: &workv1.ManifestWorkSubjectServiceAccount{
								Namespace: "default",
								Name:      "work-agent-sa",
							},
						},
					},
				},
			}

			manifestWork, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(manifestWork.Spec.ManifestConfigs[0].UpdateStrategy.Type).Should(gomega.Equal(workv1.UpdateStrategyTypeServerSideApply))
			gomega.Expect(manifestWork.Spec.ManifestConfigs[0].FeedbackRules[0].Type).Should(gomega.Equal(workv1.JSONPathsType))
			gomega.Expect(manifestWork.Spec.Executor.Subject.Type).Should(gomega.Equal(workv1.ExecutorSubjectTypeServiceAccount))
		})
	})

	ginkgo.Context("ManifestWork status validation", func() {
		ginkgo.It("should allow status updates", func() {
			work := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkName,
				},
				Spec: workv1.ManifestWorkSpec{
					Workload: workv1.ManifestsTemplate{
						Manifests: []workv1.Manifest{
							{
								RawExtension: runtime.RawExtension{
									Raw: []byte(`{"apiVersion": "v1", "kind": "ConfigMap", "metadata": {"name": "test"}}`),
								},
							},
						},
					},
				},
			}

			manifestWork, err := hubWorkClient.WorkV1().ManifestWorks(testNamespace).Create(context.TODO(), work, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			// Update status
			now := metav1.NewTime(time.Now())
			manifestWork.Status = workv1.ManifestWorkStatus{
				Conditions: []metav1.Condition{
					{
						Type:               workv1.WorkApplied,
						Status:             metav1.ConditionTrue,
						Reason:             "AppliedManifestWorkComplete",
						LastTransitionTime: now,
					},
				},
				ResourceStatus: workv1.ManifestResourceStatus{
					Manifests: []workv1.ManifestCondition{
						{
							ResourceMeta: workv1.ManifestResourceMeta{
								Ordinal:   0,
								Group:     "",
								Version:   "v1",
								Kind:      "ConfigMap",
								Resource:  "configmaps",
								Name:      "test",
								Namespace: "default",
							},
							StatusFeedbacks: workv1.StatusFeedbackResult{
								Values: []workv1.FeedbackValue{
									{
										Name: "status",
										Value: workv1.FieldValue{
											Type:   workv1.String,
											String: &[]string{"Applied"}[0],
										},
									},
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:               workv1.ManifestApplied,
									Status:             metav1.ConditionTrue,
									Reason:             "AppliedManifestComplete",
									LastTransitionTime: now,
								},
							},
						},
					},
				},
			}

			_, err = hubWorkClient.WorkV1().ManifestWorks(testNamespace).UpdateStatus(context.TODO(), manifestWork, metav1.UpdateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})
	})
})
