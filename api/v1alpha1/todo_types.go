package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TodoSpec struct {
	//+kubebuilder:validation:Required
	Title string `json:"title"`
}

// TodoStatus defines the observed state of Todo
type TodoStatus struct {
	ID int `json:"id,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced

// Todo is the Schema for the todoes API
type Todo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//+kubebuilder:validation:Required
	Spec   TodoSpec   `json:"spec"`
	Status TodoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TodoList contains a list of Todo
type TodoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Todo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Todo{}, &TodoList{})
}
