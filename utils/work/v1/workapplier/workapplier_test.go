package workapplier

import (
	"context"
	"testing"
	"time"

	v1 "k8s.io/api/apps/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clienttesting "k8s.io/client-go/testing"
	fakework "open-cluster-management.io/api/client/work/clientset/versioned/fake"
	workinformers "open-cluster-management.io/api/client/work/informers/externalversions"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

// assertActions asserts the actual actions have the expected action verb
func assertActions(t *testing.T, actualActions []clienttesting.Action, expectedVerbs ...string) {
	if len(actualActions) != len(expectedVerbs) {
		t.Fatalf("expected %d call but got: %#v", len(expectedVerbs), actualActions)
	}
	for i, expected := range expectedVerbs {
		if actualActions[i].GetVerb() != expected {
			t.Errorf("expected %s action but got: %#v", expected, actualActions[i])
		}
	}
}

// assertNoActions asserts no actions are happened
func assertNoActions(t *testing.T, actualActions []clienttesting.Action) {
	assertActions(t, actualActions)
}

func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
		},
	}
}

func newFakeWork(name, namespace string, obj runtime.Object) *workapiv1.ManifestWork {
	rawObject, _ := runtime.Encode(unstructured.UnstructuredJSONScheme, obj)

	return &workapiv1.ManifestWork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: workapiv1.ManifestWorkSpec{
			Workload: workapiv1.ManifestsTemplate{
				Manifests: []workapiv1.Manifest{
					{
						RawExtension: runtime.RawExtension{Raw: rawObject},
					},
				},
			},
			DeleteOption:    nil,
			ManifestConfigs: nil,
		},
	}
}

func TestWorkApplierWithTypedClient(t *testing.T) {
	fakeWorkClient := fakework.NewSimpleClientset()
	workInformerFactory := workinformers.NewSharedInformerFactory(fakeWorkClient, 10*time.Minute)
	fakeWorkLister := workInformerFactory.Work().V1().ManifestWorks().Lister()
	workApplier := NewWorkApplierWithTypedClient(fakeWorkClient, fakeWorkLister)

	work := newFakeWork("test", "test", newUnstructured("batch/v1", "Job", "default", "test"))
	_, err := workApplier.Apply(context.TODO(), work)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	assertActions(t, fakeWorkClient.Actions(), "create")
	if err := workInformerFactory.Work().V1().ManifestWorks().Informer().GetStore().Add(work); err != nil {
		t.Errorf("failed to add work to store with err %v", err)
	}

	// IF work is not changed, we should not update
	newWorkCopy := work.DeepCopy()
	fakeWorkClient.ClearActions()
	_, err = workApplier.Apply(context.TODO(), newWorkCopy)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	assertNoActions(t, fakeWorkClient.Actions())

	// Update work spec to update it
	newWork := newFakeWork("test", "test", newUnstructured("batch/v1", "Job", "default", "test"))
	newWork.Spec.DeleteOption = &workapiv1.DeleteOption{PropagationPolicy: workapiv1.DeletePropagationPolicyTypeOrphan}
	fakeWorkClient.ClearActions()
	appliedWork, err := workApplier.Apply(context.TODO(), newWork)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	assertActions(t, fakeWorkClient.Actions(), "patch")
	if !apiequality.Semantic.DeepEqual(appliedWork.Spec.DeleteOption, newWork.Spec.DeleteOption) {
		t.Errorf("unexpected applied work %v", appliedWork.Spec.DeleteOption)
	}
	if err := workInformerFactory.Work().V1().ManifestWorks().Informer().GetStore().Add(newWork); err != nil {
		t.Errorf("failed to add work to store with err %v", err)
	}

	// Update work annotation to update it
	newWork = appliedWork.DeepCopy()
	newWork.SetAnnotations(map[string]string{workapiv1.ManifestConfigSpecHashAnnotationKey: "hash"})
	fakeWorkClient.ClearActions()
	appliedWork, err = workApplier.Apply(context.TODO(), newWork)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	if !apiequality.Semantic.DeepEqual(appliedWork.Annotations, newWork.Annotations) {
		t.Errorf("unexpected applied work %v", appliedWork.Annotations)
	}
	assertActions(t, fakeWorkClient.Actions(), "patch")

	// Do not update if generation is not changed
	work.Spec.DeleteOption = &workapiv1.DeleteOption{PropagationPolicy: workapiv1.DeletePropagationPolicyTypeForeground}
	if err := workInformerFactory.Work().V1().ManifestWorks().Informer().GetStore().Update(work); err != nil {
		t.Errorf("failed to update work with err %v", err)
	}

	fakeWorkClient.ClearActions()
	if err := workInformerFactory.Work().V1().ManifestWorks().Informer().GetStore().Update(work); err != nil {
		t.Errorf("failed to update work with err %v", err)
	}
	_, err = workApplier.Apply(context.TODO(), newWork)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	assertNoActions(t, fakeWorkClient.Actions())

	// change generation will cause update
	work.Generation = 1
	if err := workInformerFactory.Work().V1().ManifestWorks().Informer().GetStore().Update(work); err != nil {
		t.Errorf("failed to update work with err %v", err)
	}

	fakeWorkClient.ClearActions()
	if err := workInformerFactory.Work().V1().ManifestWorks().Informer().GetStore().Update(work); err != nil {
		t.Errorf("failed to update work with err %v", err)
	}
	_, err = workApplier.Apply(context.TODO(), newWork)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	assertActions(t, fakeWorkClient.Actions(), "patch")

	fakeWorkClient.ClearActions()
	err = workApplier.Delete(context.TODO(), newWork.Namespace, newWork.Name)
	if err != nil {
		t.Errorf("failed to delete work with err %v", err)
	}
	assertActions(t, fakeWorkClient.Actions(), "delete")
}

var deploymentJson = `{
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
        "labels": {
            "app": "helloworld-agent"
        },
        "name": "helloworld-agent",
        "namespace": "default"
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "helloworld-agent"
            }
        },
		"strategy": {
            "rollingUpdate": {
                "maxSurge": "25%",
                "maxUnavailable": "25%"
            },
            "type": "RollingUpdate"
        },
        "template": {
            "metadata": {
                "labels": {
                    "app": "helloworld-agent"
                }
            },
            "spec": {
                "containers": [
                    {
                        "args": [
                            "/helloworld"
                        ],
                        "image": "quay.io/open-cluster-management/addon-examples:latest",
						"imagePullPolicy": "IfNotPresent",
                        "name": "helloworld-agent",
                        "resources": {}
                    }
                ]
            }
        }
    },
    "status":{}
}
`

// the raw in object has no creationTimestamp
func NewManifestFromJson() runtime.Object {
	obj := &unstructured.Unstructured{}
	_ = obj.UnmarshalJSON([]byte(deploymentJson))
	return obj
}

// the raw in object has creationTimestamp
func NewManifestFromDecoder() runtime.Object {
	scheme := runtime.NewScheme()
	_ = v1.AddToScheme(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
	object, _, _ := decoder.Decode([]byte(deploymentJson), nil, nil)
	return object
}

func Test_ManifestWorkEqual(t *testing.T) {
	cases := []struct {
		name         string
		requiredWork func() *workapiv1.ManifestWork
		existingWork func() *workapiv1.ManifestWork
		expected     bool
	}{
		{
			name: "required and existing without update labels",
			requiredWork: func() *workapiv1.ManifestWork {
				work := newFakeWork("test", "test", NewManifestFromJson())
				work.SetLabels(map[string]string{"test": "test"})
				work.SetAnnotations(map[string]string{"test": "test"})
				return work
			},
			existingWork: func() *workapiv1.ManifestWork {
				work := newFakeWork("test", "test", NewManifestFromDecoder())
				work.SetLabels(map[string]string{"addonname": "test"})
				work.SetLabels(map[string]string{"test": "test"})
				work.SetAnnotations(map[string]string{"test": "test"})
				return work
			},
			expected: true,
		},
		{
			name: "required and existing with different labels",
			requiredWork: func() *workapiv1.ManifestWork {
				work := newFakeWork("test", "test", NewManifestFromJson())
				work.SetLabels(map[string]string{"test": "test"})
				return work
			},
			existingWork: func() *workapiv1.ManifestWork {
				work := newFakeWork("test", "test", NewManifestFromDecoder())
				work.SetLabels(map[string]string{"addonname": "test"})
				return work
			},
			expected: false,
		},
		{
			name: "required and existing with same spec",
			requiredWork: func() *workapiv1.ManifestWork {
				work := newFakeWork("test", "test", NewManifestFromJson())
				work.SetLabels(map[string]string{"test": "test"})
				work.Spec.ManifestConfigs = []workapiv1.ManifestConfigOption{
					{
						ResourceIdentifier: workapiv1.ResourceIdentifier{},
						FeedbackRules: []workapiv1.FeedbackRule{
							{
								Type: workapiv1.WellKnownStatusType,
							},
						},
					},
				}

				return work
			},
			existingWork: func() *workapiv1.ManifestWork {
				work := newFakeWork("test", "test", NewManifestFromDecoder())
				work.SetLabels(map[string]string{"test": "test"})
				work.Spec.ManifestConfigs = []workapiv1.ManifestConfigOption{
					{
						ResourceIdentifier: workapiv1.ResourceIdentifier{},
						FeedbackRules: []workapiv1.FeedbackRule{
							{
								Type: "WellKnownStatus",
							},
						},
					},
				}
				return work
			},
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := ManifestWorkEqual(c.requiredWork(), c.existingWork())
			if c.expected != actual {
				t.Errorf("expected %v, but got %v", c.expected, actual)
			}

		})
	}
}

func TestCreateWork(t *testing.T) {
	fakeWorkClient := fakework.NewSimpleClientset()
	fakeWorkClient.ClearActions()

	workInformerFactory := workinformers.NewSharedInformerFactory(fakeWorkClient, 10*time.Minute)
	fakeWorkLister := workInformerFactory.Work().V1().ManifestWorks().Lister()
	workApplier := NewWorkApplierWithTypedClient(fakeWorkClient, fakeWorkLister)

	fakeWorkClient.PrependReactor("create", "manifestworks", func(action clienttesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, apierrors.NewAlreadyExists(workapiv1.Resource("manifestworks"), "test")
	})
	work := newFakeWork("test", "test", newUnstructured("batch/v1", "Job", "default", "test"))
	_, err := workApplier.Apply(context.TODO(), work)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	if workApplier.cache.safeToSkipApply(work, work) {
		t.Errorf("should not create work")
	}
	fakeWorkClient.ReactionChain = []clienttesting.Reactor{}
	_, err = workApplier.Apply(context.TODO(), work)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	if !workApplier.cache.safeToSkipApply(work, work) {
		t.Errorf("should create work")
	}
}
