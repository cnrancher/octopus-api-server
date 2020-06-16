package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DeviceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeviceTemplateSpec   `json:"spec,omitempty"`
	Status DeviceTemplateStatus `json:"status,omitempty"`
}

type DeviceTemplateSpec struct {
	DeviceKind          string `json:"deviceKind,omitempty"`
	DeviceVersion       string `json:"deviceVersion,omitempty"`
	DeviceGroup         string `json:"deviceGroup,omitempty"`
	DeviceResource      string `json:"deviceResource,omitempty"`
	Description         string `json:"description"`
	DefaultRevisionName string `json:"defaultRevisionName"`
}

type DeviceTemplateStatus struct {
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}
