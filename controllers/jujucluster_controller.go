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

	"cluster-api-provider-juju/pkg/juju"

	"cluster-api-provider-juju/api/v1alpha3"
	infrastructurev1alpha3 "cluster-api-provider-juju/api/v1alpha3"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	jujuCluster := &v1alpha3.JujuCluster{}
	err := r.Get(ctx, req.NamespacedName, jujuCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	switch jujuCluster.Status.State {
	case "":
		// The default case for a newly created Custom Resource
		jujuCluster.Status.State = "Pending"
		// Update the cluster object before continuing
		if err = r.Update(ctx, jujuCluster); err != nil {
			logger.Error(err, "Failed to update JujuCluster status")
		}
	case "Pending":

	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JujuClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha3.JujuCluster{}).
		Complete(r)
}
