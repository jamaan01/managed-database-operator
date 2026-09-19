package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ManagedDatabaseSpec defines the desired state of ManagedDatabase
type ManagedDatabaseSpec struct {
	// +kubebuilder:validation:MinLength=1
	Engine string `json:"engine"`

	// +kubebuilder:validation:Minimum=1
	SizeGB int32 `json:"sizeGB"`
}

// ManagedDatabaseStatus defines the observed state of ManagedDatabase.
type ManagedDatabaseStatus struct {
	CreationAttemptedAt *metav1.Time `json:"creationAttemptedAt,omitempty"`
	ExternalID          string       `json:"externalID,omitempty"`
	State               string       `json:"state,omitempty"`
	Endpoint            string       `json:"endpoint,omitempty"`
	Message             string       `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Engine",type=string,JSONPath=`.spec.engine`
// +kubebuilder:printcolumn:name="Size",type=integer,JSONPath=`.spec.sizeGB`
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`

// ManagedDatabase is the Schema for the manageddatabases API
type ManagedDatabase struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of ManagedDatabase
	// +required
	Spec ManagedDatabaseSpec `json:"spec"`

	// status defines the observed state of ManagedDatabase
	// +optional
	Status ManagedDatabaseStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ManagedDatabaseList contains a list of ManagedDatabase
type ManagedDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ManagedDatabase `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &ManagedDatabase{}, &ManagedDatabaseList{})
		return nil
	})
}
