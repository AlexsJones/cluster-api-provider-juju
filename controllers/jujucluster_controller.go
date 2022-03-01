/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/AlexsJones/cluster-api-provider-juju/api/v1alpha3"
	"github.com/AlexsJones/cluster-api-provider-juju/pkg/juju"
	"github.com/AlexsJones/cluster-api-provider-juju/pkg/utils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// JujuClusterReconciler reconciles a JujuCluster object
type JujuClusterReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	JujuClient juju.IJuju
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=jujuclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=jujuclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=jujuclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JujuCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *JujuClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get the Cluster
	jujuCluster := &v1alpha3.JujuCluster{}
	err := r.Get(ctx, req.NamespacedName, jujuCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// Fetch the JujuConfiguration object
	jujuConfiguration, err := utils.FetchJujuConfigurationObject(jujuCluster, r.Client, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	finalizerName := "infrastructure.cluster.x-k8s.io/finalizer"
	// Check to see if the Cluster has a finalizer
	if jujuCluster.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !ctrlutil.ContainsFinalizer(jujuCluster, finalizerName) {
			ctrlutil.AddFinalizer(jujuCluster, finalizerName)
			if err := r.Update(ctx, jujuCluster); err != nil {
				return ctrl.Result{}, err
			}
		}
		// Configuration object finalizer
		if !ctrlutil.ContainsFinalizer(jujuConfiguration, finalizerName) {
			ctrlutil.AddFinalizer(jujuConfiguration, finalizerName)
			if err := r.Update(ctx, jujuConfiguration); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// Check to see if the Cluster is under deletion
		if controllerutil.ContainsFinalizer(jujuCluster, finalizerName) {

			// The object is being deleted
			jujuCluster.Status.State = "Deleting"
			if err := r.Status().Update(ctx, jujuCluster); err != nil {
				return ctrl.Result{}, err
			}

			jujuCluster = &v1alpha3.JujuCluster{}
			err := r.Get(ctx, req.NamespacedName, jujuCluster)
			if err != nil {
				if apierrors.IsNotFound(err) {
					return ctrl.Result{}, nil
				}
				return ctrl.Result{}, err
			}

			if err := r.JujuClient.DestroyCluster(jujuConfiguration, jujuCluster); err != nil {
				return ctrl.Result{}, err
			}

			// remove our finalizer from the Cluster and update it.
			controllerutil.RemoveFinalizer(jujuCluster, finalizerName)
			if err := r.Update(ctx, jujuCluster); err != nil {
				return ctrl.Result{}, err
			}

			// remove our finalizer from the Config and update it.
			controllerutil.RemoveFinalizer(jujuConfiguration, finalizerName)
			if err := r.Update(ctx, jujuConfiguration); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	switch jujuCluster.Status.State {
	case "":
		status, err := r.JujuClient.GetClusterStatus(jujuConfiguration)
		if err != nil {
			return ctrl.Result{}, err
		}
		switch status {
		case juju.E_JUJU_CLUSTER_STATUS_RUNNING:
			jujuCluster.Status.State = "Running"
			return ctrl.Result{}, r.Status().Update(ctx, jujuCluster)
		}
		// Check if the cluster exists, otherwise create it
		if err := r.JujuClient.CreateCluster(jujuConfiguration, jujuCluster); err != nil {
			return ctrl.Result{}, err
		}

		jujuCluster.Status.State = "Provisioning"
		// Update the cluster object before continuing
		if err = r.Status().Update(ctx, jujuCluster); err != nil {
			logger.Error(err, "Failed to update JujuCluster status")
		}
	case "Unknown":
		fallthrough
	case "Provisioning":
		status, err := r.JujuClient.GetClusterStatus(jujuConfiguration)
		if err != nil {
			return ctrl.Result{}, err
		}
		switch status {
		case juju.E_JUJU_CLUSTER_STATUS_RUNNING:
			jujuCluster.Status.State = "Running"

		case juju.E_JUJU_CLUSTER_STATUS_UNKNOWN:
			jujuCluster.Status.State = "Provisioning"
		}
		if err = r.Status().Update(ctx, jujuCluster); err != nil {
			logger.Error(err, "Failed to update JujuCluster status")
		}
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *JujuClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha3.JujuCluster{}).
		Complete(r)
}
