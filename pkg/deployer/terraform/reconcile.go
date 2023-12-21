package terraform

import (
	"context"
	"time"

	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r JobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	job := &batchv1.Job{}

	err := r.KubeClient.Get(ctx, req.NamespacedName, job)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	err = r.syncApplicationRevisionStatus(ctx, job)
	if err != nil && !model.IsNotFound(err) {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r JobReconciler) Setup(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Job{}).
		Complete(r)
}

// syncApplicationRevisionStatus sync the application revision status.
func (r JobReconciler) syncApplicationRevisionStatus(ctx context.Context, job *batchv1.Job) (err error) {
	appRevisionID, ok := job.Labels[_resourceRevisionIDLabel]
	if !ok {
		// Not a deployer job.
		return nil
	}

	appRevision, err := r.ModelClient.ResourceRevisions().Get(ctx, object.ID(appRevisionID))
	if err != nil {
		return err
	}

	// If the application revision status is not running, then skip it.
	if !status.ResourceRevisionStatusReady.IsUnknown(appRevision) {
		return nil
	}

	if job.Status.Succeeded == 0 && job.Status.Failed == 0 {
		return nil
	}

	status.ResourceRevisionStatusReady.True(appRevision, "")

	if job.Status.Succeeded > 0 {
		r.Logger.Info("succeed", "application-revision", appRevisionID)
	}

	if job.Status.Failed > 0 {
		r.Logger.Info("failed", "application-revision", appRevisionID)
		status.ResourceRevisionStatusReady.False(appRevision, "")
	}

	// Get job pods logs.
	record, err := r.getJobPodsLogs(ctx, job.Name)
	if err != nil {
		r.Logger.Error(err, "failed to get job pod logs", "application-revision", appRevisionID)
		record = err.Error()
	}

	// Report to application revision.
	appRevision.Record = record
	appRevision.Status.SetSummary(status.WalkResourceRevision(&appRevision.Status))
	appRevision.Duration = int(time.Since(*appRevision.CreateTime).Seconds())

	appRevision, err = r.ModelClient.ResourceRevisions().UpdateOne(appRevision).
		SetStatus(appRevision.Status).
		SetRecord(appRevision.Record).
		SetDuration(appRevision.Duration).
		Save(ctx)
	if err != nil {
		return err
	}

	return revisionbus.Notify(ctx, r.ModelClient, appRevision)
}

// getJobPodsLogs returns the logs of all pods of a job.
func (r JobReconciler) getJobPodsLogs(ctx context.Context, jobName string) (string, error) {
	clientSet, err := kubernetes.NewForConfig(r.Kubeconfig)
	if err != nil {
		return "", err
	}
	ls := "job-name=" + jobName

	pods, err := clientSet.CoreV1().Pods(types.WalrusSystemNamespace).
		List(ctx, metav1.ListOptions{LabelSelector: ls})
	if err != nil {
		return "", err
	}

	var logs string

	for _, pod := range pods.Items {
		var podLogs []byte

		podLogs, err = clientSet.CoreV1().Pods(types.WalrusSystemNamespace).
			GetLogs(pod.Name, &corev1.PodLogOptions{}).
			DoRaw(ctx)
		if err != nil {
			return "", err
		}
		logs += string(podLogs)
	}

	return logs, nil
}
