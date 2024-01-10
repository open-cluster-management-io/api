package client

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubetypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog/v2"

	workv1client "open-cluster-management.io/api/client/work/clientset/versioned/typed/work/v1"
	workv1lister "open-cluster-management.io/api/client/work/listers/work/v1"
	"open-cluster-management.io/api/cloudevents/generic"
	"open-cluster-management.io/api/cloudevents/generic/types"
	"open-cluster-management.io/api/cloudevents/work/source/codec"
	"open-cluster-management.io/api/cloudevents/work/utils"
	"open-cluster-management.io/api/cloudevents/work/watcher"
	workv1 "open-cluster-management.io/api/work/v1"
)

// ManifestWorkSourceClient implements the ManifestWorkInterface.
type ManifestWorkSourceClient struct {
	cloudEventsClient *generic.CloudEventSourceClient[*workv1.ManifestWork]
	watcher           *watcher.ManifestWorkWatcher
	lister            workv1lister.ManifestWorkLister
	namespace         string
}

var manifestWorkGR = schema.GroupResource{Group: workv1.GroupName, Resource: "manifestworks"}

var _ workv1client.ManifestWorkInterface = &ManifestWorkSourceClient{}

func NewManifestWorkSourceClient(cloudEventsClient *generic.CloudEventSourceClient[*workv1.ManifestWork], watcher *watcher.ManifestWorkWatcher) *ManifestWorkSourceClient {
	return &ManifestWorkSourceClient{
		cloudEventsClient: cloudEventsClient,
		watcher:           watcher,
	}
}

func (c *ManifestWorkSourceClient) SetLister(lister workv1lister.ManifestWorkLister) {
	c.lister = lister
}

func (mw *ManifestWorkSourceClient) SetNamespace(namespace string) {
	mw.namespace = namespace
}

func (c *ManifestWorkSourceClient) Create(ctx context.Context, manifestWork *workv1.ManifestWork, opts metav1.CreateOptions) (*workv1.ManifestWork, error) {
	_, err := c.lister.ManifestWorks(c.namespace).Get(manifestWork.Name)
	if err == nil {
		return nil, errors.NewAlreadyExists(manifestWorkGR, manifestWork.Name)
	}

	if errors.IsNotFound(err) {
		eventDataType, err := getWorkCloudEventDataType(manifestWork)
		if err != nil {
			return nil, err
		}

		eventType := types.CloudEventsType{
			CloudEventsDataType: *eventDataType,
			SubResource:         types.SubResourceSpec,
		}

		newWork := manifestWork.DeepCopy()
		if err := c.cloudEventsClient.Publish(ctx, eventType, newWork); err != nil {
			return nil, err
		}

		// add the new work to the ManifestWorkInformer local cache.
		c.watcher.Receive(watch.Event{Type: watch.Added, Object: newWork})
		return newWork.DeepCopy(), nil
	}

	return nil, err
}

func (c *ManifestWorkSourceClient) Update(ctx context.Context, manifestWork *workv1.ManifestWork, opts metav1.UpdateOptions) (*workv1.ManifestWork, error) {
	return nil, errors.NewMethodNotSupported(manifestWorkGR, "update")
}

func (c *ManifestWorkSourceClient) UpdateStatus(ctx context.Context, manifestWork *workv1.ManifestWork, opts metav1.UpdateOptions) (*workv1.ManifestWork, error) {
	return nil, errors.NewMethodNotSupported(manifestWorkGR, "updatestatus")
}

func (c *ManifestWorkSourceClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	work, err := c.lister.ManifestWorks(c.namespace).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	eventDataType, err := getWorkCloudEventDataType(work)
	if err != nil {
		return err
	}

	eventType := types.CloudEventsType{
		CloudEventsDataType: *eventDataType,
		SubResource:         types.SubResourceSpec,
	}

	deletingWork := work.DeepCopy()
	now := metav1.Now()
	deletingWork.DeletionTimestamp = &now

	if err := c.cloudEventsClient.Publish(ctx, eventType, deletingWork); err != nil {
		return err
	}

	// update the deleting work in the ManifestWorkInformer local cache.
	c.watcher.Receive(watch.Event{Type: watch.Modified, Object: deletingWork})
	return nil
}

func (c *ManifestWorkSourceClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return errors.NewMethodNotSupported(manifestWorkGR, "deletecollection")
}

func (c *ManifestWorkSourceClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*workv1.ManifestWork, error) {
	klog.V(4).Infof("getting manifestwork %s", name)
	return c.lister.ManifestWorks(c.namespace).Get(name)
}

func (c *ManifestWorkSourceClient) List(ctx context.Context, opts metav1.ListOptions) (*workv1.ManifestWorkList, error) {
	klog.V(4).Infof("list manifestworks")
	return &workv1.ManifestWorkList{}, nil
}

func (c *ManifestWorkSourceClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	// TODO (skeeey) consider resync the manifestworks when the ManifestWorkInformer reconnected
	return c.watcher, nil
}

func (c *ManifestWorkSourceClient) Patch(ctx context.Context, name string, pt kubetypes.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *workv1.ManifestWork, err error) {
	klog.V(4).Infof("patching manifestwork %s", name)

	if len(subresources) != 0 {
		return nil, fmt.Errorf("unsupported to update subresources %v", subresources)
	}

	lastWork, err := c.lister.ManifestWorks(c.namespace).Get(name)
	if err != nil {
		return nil, err
	}

	patchedWork, err := utils.Patch(pt, lastWork, data)
	if err != nil {
		return nil, err
	}

	eventDataType, err := getWorkCloudEventDataType(patchedWork)
	if err != nil {
		return nil, err
	}

	eventType := types.CloudEventsType{
		CloudEventsDataType: *eventDataType,
		SubResource:         types.SubResourceSpec,
	}

	newWork := patchedWork.DeepCopy()
	if err := c.cloudEventsClient.Publish(ctx, eventType, newWork); err != nil {
		return nil, err
	}

	// refresh the work in the ManifestWorkInformer local cache with patched work.
	c.watcher.Receive(watch.Event{Type: watch.Modified, Object: newWork})
	return newWork.DeepCopy(), nil
}

func getWorkCloudEventDataType(work *workv1.ManifestWork) (*types.CloudEventsDataType, error) {
	eventDataType, ok := work.Annotations[codec.CloudEventsDataTypeAnnotationKey]
	if !ok {
		return &types.CloudEventsDataType{
			Group:    "io.open-cluster-management.works",
			Version:  "v1alpha1",
			Resource: "manifestbundles",
		}, nil
	}

	return types.ParseCloudEventsDataType(eventDataType)
}
