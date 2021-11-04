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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type CartItem struct {
	Name  string `json:"name"`
	Count int32  `json:"count"`
}

// CartSpec defines the desired state of Cart
type CartSpec struct {
	// a list of items
	CartItems []CartItem `json:"cartitems,omitempty"`
	//+kubebuilder:default:=false
	//+kubebuilder:validation:Optional
	CheckOut bool `json:"checkout"`
}

// CartStatus defines the observed state of Cart
type CartStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
	Total   int32  `json:"total,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cart is the Schema for the carts API
type Cart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CartSpec   `json:"spec,omitempty"`
	Status CartStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CartList contains a list of Cart
type CartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cart{}, &CartList{})
}
