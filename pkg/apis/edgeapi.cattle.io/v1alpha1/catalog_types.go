package v1alpha1

import (
	"github.com/rancher/wrangler/pkg/condition"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Catalog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CatalogSpec   `json:"spec,omitempty"`
	Status CatalogStatus `json:"status,omitempty"`
}

type CatalogSpec struct {
	URL       string     `json:"url"`
	Username  string     `json:"username,omitempty"`
	Password  string     `json:"password,omitempty"`
	IndexFile *IndexFile `json:"indexFile" yaml:"indexFile"`
}

type CatalogStatus struct {
	LastRefreshTimestamp string             `json:"lastRefreshTimestamp,omitempty"`
	Conditions           []CatalogCondition `json:"conditions,omitempty"`
}

type IndexFile struct {
	Entries map[string]ChartVersions `json:"entries" yaml:"entries"`
}

type ChartVersions []*ChartVersion

type ChartVersion struct {
	ChartMetadata `yaml:",inline"`
	URLs          []string `json:"urls" yaml:"urls"`
	Digest        string   `json:"digest,omitempty" yaml:"digest,omitempty"`
}

type ChartMetadata struct {
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	KubeVersion string `json:"kubeVersion,omitempty" yaml:"kubeVersion,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
}

var (
	CatalogConditionCreated   condition.Cond = "Created"
	CatalogConditionRefreshed condition.Cond = "Refreshed"
	CatalogConditionProcessed condition.Cond = "Processed"
)

type CatalogCondition struct {
	Type CatalogConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type CatalogConditionType string
