package workapplier

import (
	"context"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
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
	newWork, err = workApplier.Apply(context.TODO(), newWork)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
	}
	assertActions(t, fakeWorkClient.Actions(), "patch")
	if err := workInformerFactory.Work().V1().ManifestWorks().Informer().GetStore().Add(newWork); err != nil {
		t.Errorf("failed to add work to store with err %v", err)
	}

	// Update work annotation to update it
	newWork = newWork.DeepCopy()
	newWork.SetAnnotations(map[string]string{workapiv1.ManifestConfigSpecHashAnnotationKey: "hash"})
	fakeWorkClient.ClearActions()
	_, err = workApplier.Apply(context.TODO(), newWork)
	if err != nil {
		t.Errorf("failed to apply work with err %v", err)
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
