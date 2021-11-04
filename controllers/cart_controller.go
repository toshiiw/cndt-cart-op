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

// CartReconciler reconciles a Cart object
type CartReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cart.example.valinux.co.jp,resources=carts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cart.example.valinux.co.jp,resources=carts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cart.example.valinux.co.jp,resources=carts/finalizers,verbs=update
//+kubebuilder:rbac:groups=cart.example.valinux.co.jp,resources=pricelists,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cart object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *CartReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("cart", req.NamespacedName)

	var cart cartv1alpha1.Cart
	if err := r.Get(ctx, req.NamespacedName, &cart); err != nil {
		log.Error(err, "unable to fetch")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	m := make(map[string]int32)
	for _, i := range cart.Spec.CartItems {
		m[i.Name] = i.Count
	}
	log.Info("Processing", "items", m)
	var prices cartv1alpha1.PriceListList
	if err := r.List(ctx, &prices, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "cannot list prices")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var total int32
	for _, pp := range prices.Items {
		for _, p := range pp.Spec.Prices {
			count, found := m[p.Name]
			if found {
				total += count * p.Price
				delete(m, p.Name)
			}
		}
	}
	log.Info("calculated total", "price", total)
	cart.Status.Total = total
	if len(m) > 0 {
		cart.Status.State = "Error"
		cart.Status.Message = "Unknown objects in cart"
	} else {
		cart.Status.State = "Ready"
		cart.Status.Message = ""
	}
	if cart.Spec.CheckOut {
		if cart.Status.State == "Error" {
			cart.Spec.CheckOut = false
			if err := r.Update(ctx, &cart); err != nil {
				log.Error(err, "unable to clear checkout")
			} else {
				return ctrl.Result{}, nil
			}
		} else if cart.Status.State == "Ready" {
			cart.Status.State = "CheckedOut"
		}
	}

	err := r.Status().Update(ctx, &cart)
	if err != nil {
		log.Error(err, "unable to update")
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *CartReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cartv1alpha1.Cart{}).
		Complete(r)
}
