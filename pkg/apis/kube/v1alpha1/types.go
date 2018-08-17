package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DatabaseList specifies the list of custom Database resources
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Database `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Database specifies the singular database custom resource
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              DatabaseSpec   `json:"spec"`
	Status            DatabaseStatus `json:"status,omitempty"`
}

// DatabaseSpec specifies the fields required in the spec
type DatabaseSpec struct {
	Name string `json:"name"`
}

// DatabaseStatus tracks the lifecycle of the database creation and deletion process
type DatabaseStatus struct {
	Status   string `json:"status"`
	Password string `json:"password"`
	Username string `json:"username"`
}
