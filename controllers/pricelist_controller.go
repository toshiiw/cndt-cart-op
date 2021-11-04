/*
Copyright 2021.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cartv1alpha1 "github.com/toshiiw/cndt-cart-op/api/v1alpha1"
)

// PriceListReconciler reconciles a PriceList object
type PriceListReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cart.example.valinux.co.jp,resources=pricelists,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cart.example.valinux.co.jp,resources=pricelists/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cart.example.valinux.co.jp,resources=pricelists/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PriceList object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *PriceListReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var carts cartv1alpha1.CartList
	if err := r.List(ctx, &carts, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "cannot list carts")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	var lerr error
	for _, c := range carts.Items {
		if c.Status.State == "CheckedOut" {
			continue
		}
		c.Status.State = "NeedsUpdate"
		log.Info("needsupdate", "cart", c)
		err := r.Status().Update(ctx, &c)
		if err != nil {
			log.Error(err, "unable to update cart.status")
			lerr = err
		}
	}

	return ctrl.Result{}, lerr
}

// SetupWithManager sets up the controller with the Manager.
func (r *PriceListReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cartv1alpha1.PriceList{}).
		Complete(r)
}
