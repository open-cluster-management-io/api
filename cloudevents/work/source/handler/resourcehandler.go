package handler

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog/v2"
	workv1lister "open-cluster-management.io/api/client/work/listers/work/v1"

	"open-cluster-management.io/api/cloudevents/generic"
	"open-cluster-management.io/api/cloudevents/generic/types"
	"open-cluster-management.io/api/cloudevents/work/watcher"
	workv1 "open-cluster-management.io/api/work/v1"
)

const ManifestWorkFinalizer = "cloudevents.open-cluster-management.io/manifest-work-cleanup"
const ManifestsDeleted = "Deleted"

// NewManifestWorkSourceHandler returns a ResourceHandler for a ManifestWork source client. It sends the kube events
// with ManifestWorWatcher after CloudEventSourceClient received the ManifestWork status from agent, then the
// ManifestWorkInformer handles the kube events in its local cache.
func NewManifestWorkSourceHandler(lister workv1lister.ManifestWorkLister, watcher *watcher.ManifestWorkWatcher) generic.ResourceHandler[*workv1.ManifestWork] {
	return func(action types.ResourceAction, work *workv1.ManifestWork) error {
		switch action {
		case types.StatusModified:
			lastWork, err := lister.ManifestWorks(work.Namespace).Get(work.Name)
			if err != nil {
				return err
			}

			if work.Generation < lastWork.Generation {
				klog.Infof("The work %s generation %d is less than cached generation %d, ignore",
					work.UID, work.Generation, lastWork.Generation)
				return nil
			}

			// no status change
			if equality.Semantic.DeepEqual(lastWork.Status, work.Status) {
				return nil
			}

			// restore the fields that are maintained by local agent
			work.Labels = lastWork.Labels
			work.Annotations = lastWork.Annotations
			work.DeletionTimestamp = lastWork.DeletionTimestamp
			work.Spec = lastWork.Spec

			if meta.IsStatusConditionTrue(work.Status.Conditions, ManifestsDeleted) {
				work.Finalizers = []string{}
				watcher.Receive(watch.Event{Type: watch.Deleted, Object: work})
				return nil
			}

			// the work is handled by agent, we make sure a finalizer here
			work.Finalizers = mergeFinalizers(lastWork.Finalizers)
			watcher.Receive(watch.Event{Type: watch.Modified, Object: work})
		default:
			return fmt.Errorf("unsupported resource action %s", action)
		}

		return nil
	}
}

func mergeFinalizers(workFinalizers []string) []string {
	has := false
	for _, f := range workFinalizers {
		if f == ManifestWorkFinalizer {
			has = true
			break
		}
	}

	if !has {
		workFinalizers = append(workFinalizers, ManifestWorkFinalizer)
	}

	return workFinalizers
}
