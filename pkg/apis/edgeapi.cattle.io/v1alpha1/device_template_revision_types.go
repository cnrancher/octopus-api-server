package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DeviceTemplateRevision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeviceTemplateRevisionSpec   `json:"spec,omitempty"`
	Status DeviceTemplateRevisionStatus `json:"status,omitempty"`
}

type DeviceTemplateRevisionSpec struct {
	DisplayName              string            `json:"displayName"`
	Enabled                  *bool             `json:"enabled,omitempty"`
	DeviceTemplateName       string            `json:"deviceTemplateName"`
	DeviceTemplateAPIVersion string            `json:"deviceTemplateAPIVersion"`
	Labels                   map[string]string `json:"labels,omitempty"`

	TemplateSpec *runtime.RawExtension `json:"templateSpec,omitempty"`
}

type DeviceTemplateRevisionStatus struct {
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}
