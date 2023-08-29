package client

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog/v2"

	workv1client "open-cluster-management.io/api/client/work/clientset/versioned/typed/work/v1"
	workv1lister "open-cluster-management.io/api/client/work/listers/work/v1"
	"open-cluster-management.io/api/cloudevents/generic"
	"open-cluster-management.io/api/cloudevents/generic/types"
	"open-cluster-management.io/api/cloudevents/work/agent/codec"
	"open-cluster-management.io/api/cloudevents/work/watcher"
	workv1 "open-cluster-management.io/api/work/v1"
)

const ManifestsDeleted = "Deleted"

const (
	UpdateRequestAction = "update_request"
	DeleteRequestAction = "delete_request"
)

// ManifestWorksAgentClient implements the ManifestWorkInterface. It sends the manifestworks status back to source by
// CloudEventAgentClient.
type ManifestWorksAgentClient struct {
	cloudEventsClient generic.CloudEventAgentClient[*workv1.ManifestWork]
	watcher           *watcher.ManifestWorkWatcher
	lister            workv1lister.ManifestWorkNamespaceLister
}

var _ workv1client.ManifestWorkInterface = &ManifestWorksAgentClient{}

func (c *ManifestWorksAgentClient) Create(ctx context.Context, manifestWork *workv1.ManifestWork, opts metav1.CreateOptions) (*workv1.ManifestWork, error) {
	klog.Fatal("Create function for ManifestWorksAgentClient is unsupported")
	return nil, nil
}

func (c *ManifestWorksAgentClient) Update(ctx context.Context, manifestWork *workv1.ManifestWork, opts metav1.UpdateOptions) (*workv1.ManifestWork, error) {
	// TODO (skeeey) using patch instead
	klog.V(4).Infof("updating manifestwork %s", manifestWork.Name)

	if !manifestWork.DeletionTimestamp.IsZero() && len(manifestWork.Finalizers) == 0 {
		// the finalizers of a deleting manifestwork are removed on a managed cluster, marking the manifest work status
		// to deleted and send it back to source
		meta.SetStatusCondition(&manifestWork.Status.Conditions, metav1.Condition{
			Type:    ManifestsDeleted,
			Status:  metav1.ConditionTrue,
			Reason:  "ManifestsDeleted",
			Message: fmt.Sprintf("The manifests are deleted from the cluster %s", manifestWork.Namespace),
		})

		eventDataType, err := types.ParseCloudEventsDataType(manifestWork.Annotations[codec.CloudEventsDataTypeAnnotationKey])
		if err != nil {
			return nil, err
		}

		eventType := types.CloudEventsType{
			CloudEventsDataType: *eventDataType,
			SubResource:         types.SubResourceStatus,
			Action:              DeleteRequestAction,
		}

		if err := c.cloudEventsClient.Publish(ctx, eventType, manifestWork); err != nil {
			return nil, err
		}

		// also send the deleted event to delete the manifestwork from the ManifestWorkInformer
		c.watcher.Receive(watch.Event{Type: watch.Deleted, Object: manifestWork})
		return manifestWork, nil
	}

	return nil, nil
}

func (c *ManifestWorksAgentClient) UpdateStatus(ctx context.Context, manifestWork *workv1.ManifestWork, opts metav1.UpdateOptions) (*workv1.ManifestWork, error) {
	// TODO (skeeey) using patch instead
	klog.V(4).Infof("updating manifestwork %s status", manifestWork.Name)

	_, err := c.lister.Get(manifestWork.Name)
	if err != nil {
		return nil, err
	}

	eventDataType, err := types.ParseCloudEventsDataType(manifestWork.Annotations[codec.CloudEventsDataTypeAnnotationKey])
	if err != nil {
		return nil, err
	}

	updatedWork := manifestWork.DeepCopy()

	eventType := types.CloudEventsType{
		CloudEventsDataType: *eventDataType,
		SubResource:         types.SubResourceStatus,
		Action:              UpdateRequestAction,
	}

	if err := c.cloudEventsClient.Publish(ctx, eventType, updatedWork); err != nil {
		return nil, err
	}
	return updatedWork, nil
}

func (c *ManifestWorksAgentClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	klog.Fatal("Delete function for ManifestWorksAgentClient is unsupported")
	return nil
}

func (c *ManifestWorksAgentClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	klog.Fatal("DeleteCollection function for ManifestWorksAgentClient is unsupported")
	return nil
}

func (c *ManifestWorksAgentClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*workv1.ManifestWork, error) {
	klog.V(4).Infof("getting manifestwork %s", name)
	return c.lister.Get(name)
}

func (c *ManifestWorksAgentClient) List(ctx context.Context, opts metav1.ListOptions) (*workv1.ManifestWorkList, error) {
	klog.V(4).Infof("sync manifestworks")
	// send resync request to fetch manifestworks from source when the ManifestWorkInformer status
	if err := c.cloudEventsClient.Resync(ctx); err != nil {
		return nil, err
	}

	return &workv1.ManifestWorkList{}, nil
}

func (mw *ManifestWorksAgentClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	// TODO (skeeey) consider resync the manifestworks when the ManifestWorkInformer reconnected
	return mw.watcher, nil
}

func (mw *ManifestWorksAgentClient) Patch(ctx context.Context, name string, pt kubetypes.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *workv1.ManifestWork, err error) {
	klog.Fatal("Patch function for ManifestWorksAgentClient has not been implemented")
	return nil, nil
}
