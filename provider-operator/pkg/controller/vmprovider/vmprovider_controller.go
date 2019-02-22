package vmprovider

import (
	"context"
	"encoding/base64"

	kubevirtv1alpha1 "github.com/pkliczewski/provider-operators/provider-operator/pkg/apis/kubevirt/v1alpha1"
	vmclient "github.com/pkliczewski/provider-operators/provider-operator/pkg/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_vmprovider")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new VmProvider Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileVmProvider{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("vmprovider-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource VmProvider
	err = c.Watch(&source.Kind{Type: &kubevirtv1alpha1.VmProvider{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner VmProvider
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kubevirtv1alpha1.VmProvider{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileVmProvider{}

// ReconcileVmProvider reconciles a VmProvider object
type ReconcileVmProvider struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a VmProvider object and makes changes based on the state read
// and what is in the VmProvider.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileVmProvider) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling VmProvider")

	// Fetch the VmProvider instance
	instance := &kubevirtv1alpha1.VmProvider{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// VmProvider already validated - don't requeue
	if instance.Status.Validated == true {
		reqLogger.Info("Skip reconcile: Provider already validated", "VmProvider.Namespace", instance.Namespace, "VmProvider.Name", instance.Name)
		return reconcile.Result{}, nil
	}

	err = validate(instance)
	if err != nil {
		instance.Status.Validated = false
		instance.Status.Message = err.Error()
	} else {
		instance.Status.Validated = true
	}

	reqLogger.Info("Provider validation done", "VmProvider.Namespace", instance.Namespace, "VmProvider.Name", instance.Name, "VmProvider.Status.Validated", instance.Status.Validated, "Message", instance.Status.Message)

	err = r.client.Status().Update(context.TODO(), instance)
	reqLogger.Info("Provider Updated", "VmProvider.Namespace", instance.Namespace, "VmProvider.Name", instance.Name, "VmProvider.Status.Validated", instance.Status.Validated)
	if err != nil {
		// update failed - requeue the request
		reqLogger.Info("Provider failed", err.Error())
		return reconcile.Result{}, err
	}

	// done - don't requeue
	return reconcile.Result{}, nil

	// // Define a new Pod object
	// pod := newPodForCR(instance)

	// // Set VmProvider instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Check if this Pod already exists
	// found := &corev1.Pod{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	// 	err = r.client.Create(context.TODO(), pod)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Pod created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Pod already exists - don't requeue
	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	// return reconcile.Result{}, nil
}

func validate(cr *kubevirtv1alpha1.VmProvider) error {
	ctx := context.Background()

	pass, err := base64.StdEncoding.DecodeString(cr.Spec.Password)
	if err != nil {
		return err
	}

	c, err := vmclient.NewClient(ctx, cr.Spec.Url, cr.Spec.Username, string(pass))
	if err != nil {
		return err
	}
	defer c.Logout(ctx)
	c.GetVMs(ctx)
	return nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *kubevirtv1alpha1.VmProvider) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
