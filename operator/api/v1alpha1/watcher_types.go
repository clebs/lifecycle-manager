/*
Copyright 2022.

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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ManagedBylabel = "operator.kyma-project.io/managed-by"

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WatcherSpec defines the desired state of Watcher.
type WatcherSpec struct {
	// ServiceInfo describes the service information of the operator
	ServiceInfo Service `json:"serviceInfo"`

	// LabelsToWatch describes the labels that should be watched
	LabelsToWatch map[string]string `json:"labelsToWatch"`

	// Field describes the subresource that should be watched
	// Value can be one of ("spec", "status")
	Field FieldName `json:"field"`
}

// +kubebuilder:validation:Enum=spec;status;
type FieldName string

const (
	// SpecField represents FieldName spec, which indicates that resource spec will be watched.
	SpecField FieldName = "spec"
	// StatusField represents FieldName status, which indicates that only resource status will be watched.
	StatusField FieldName = "status"
)

// Service describes the service specification for the corresponding operator container.
type Service struct {
	// Port describes the service port.
	Port int64 `json:"port"`

	// Name describes the service name.
	Name string `json:"name"`

	// Namespace describes the service namespace.
	Namespace string `json:"namespace"`
}

// +kubebuilder:validation:Enum=Processing;Deleting;Ready;Error
type WatcherState string

// Valid Watcher States.
const (
	// WatcherStateReady signifies Watcher is ready and has been installed successfully.
	WatcherStateReady WatcherState = "Ready"

	// WatcherStateProcessing signifies Watcher is reconciling and is in the process of installation.
	WatcherStateProcessing WatcherState = "Processing"

	// WatcherStateError signifies an error for Watcher. This signifies that the Installation
	// process encountered an error.
	WatcherStateError WatcherState = "Error"

	// WatcherStateDeleting signifies Watcher is being deleted.
	WatcherStateDeleting WatcherState = "Deleting"
)

// WatcherStatus defines the observed state of Watcher.
type WatcherStatus struct {
	// State signifies current state of a Watcher.
	// Value can be one of ("Ready", "Processing", "Error", "Deleting")
	State WatcherState `json:"state"`

	// List of status conditions to indicate the status of a Watcher.
	// +kubebuilder:validation:Optional
	Conditions []WatcherCondition `json:"conditions"`

	// ObservedGeneration
	// +kubebuilder:validation:Optional
	ObservedGeneration int64 `json:"observedGeneration"`
}

// WatcherCondition describes condition information for Watcher.
type WatcherCondition struct {
	// Type is used to reflect what type of condition we are dealing with.
	// Most commonly WatcherConditionTypeReady it is used as extension marker in the future.
	Type WatcherConditionType `json:"type"`

	// Status of the Watcher Condition.
	// Value can be one of ("True", "False", "Unknown").
	Status WatcherConditionStatus `json:"status"`

	// Human-readable message indicating details about the last status transition.
	// +kubebuilder:validation:Optional
	Message string `json:"message"`

	// Machine-readable text indicating the reason for the condition's last transition.
	// +kubebuilder:validation:Optional
	Reason string `json:"reason"`

	// Timestamp for when Watcher last transitioned from one status to another.
	// +kubebuilder:validation:Optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime"`
}

// +kubebuilder:validation:Enum=Ready
type WatcherConditionType string

const (
	// WatcherConditionTypeReady represents WatcherConditionType Ready,
	// meaning as soon as its true we will reconcile Watcher into WatcherStateReady.
	WatcherConditionTypeReady WatcherConditionType = "Ready"
)

// +kubebuilder:validation:Enum=True;False;Unknown;
type WatcherConditionStatus string

// Valid WatcherConditionStatus.
const (
	// ConditionStatusTrue signifies WatcherConditionStatus true.
	ConditionStatusTrue WatcherConditionStatus = "True"

	// ConditionStatusFalse signifies WatcherConditionStatus false.
	ConditionStatusFalse WatcherConditionStatus = "False"

	// ConditionStatusUnknown signifies WatcherConditionStatus unknown.
	ConditionStatusUnknown WatcherConditionStatus = "Unknown"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="State",type=string,JSONPath=".status.state"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Watcher is the Schema for the watchers API.
type Watcher struct {
	metav1.TypeMeta `json:",inline"`

	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata"`

	// +kubebuilder:validation:Optional
	Spec WatcherSpec `json:"spec"`

	// +kubebuilder:validation:Optional
	Status WatcherStatus `json:"status"`
}

func (w *Watcher) SetObservedGeneration() *Watcher {
	w.Status.ObservedGeneration = w.Generation
	return w
}

func (w *Watcher) GetModuleName() string {
	if w.Labels == nil {
		return ""
	}
	return w.Labels[ManagedBylabel]
}

func (w *Watcher) AddOrUpdateReadyCondition(state WatcherConditionStatus, msg string) {
	lastTransitionTime := &metav1.Time{Time: time.Now()}
	if len(w.Status.Conditions) == 0 {
		w.Status.Conditions = []WatcherCondition{{
			Type:               WatcherConditionTypeReady,
			Status:             state,
			Message:            msg,
			LastTransitionTime: lastTransitionTime,
		}}
	}
	for i := range w.Status.Conditions {
		condition := &w.Status.Conditions[i]
		if condition.Type == WatcherConditionTypeReady {
			condition.Status = state
			condition.LastTransitionTime = lastTransitionTime
		}
	}
}

//+kubebuilder:object:root=true

// WatcherList contains a list of Watcher.
type WatcherList struct {
	metav1.TypeMeta `json:",inline"`

	// +kubebuilder:validation:Optional
	metav1.ListMeta `json:"metadata"`
	Items           []Watcher `json:"items"`
}

func init() { //nolint:gochecknoinits
	SchemeBuilder.Register(&Watcher{}, &WatcherList{})
}
