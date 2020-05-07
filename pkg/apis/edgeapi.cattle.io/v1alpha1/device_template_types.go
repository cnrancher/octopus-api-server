package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	DeviceKind     string                `json:"deviceKind,omitempty"`
	DeviceVersion  string                `json:"deviceVersion,omitempty"`
	DeviceGroup    string                `json:"deviceGroup,omitempty"`
	DeviceResource string                `json:"deviceResource,omitempty"`
	Labels         map[string]string     `json:"labels,omitempty"`
	TemplateSpec   *runtime.RawExtension `json:"templateSpec,omitempty"`
}

type DeviceTemplateStatus struct {
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}
