/*
 * Kubernetes
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: v1.10.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

import (
	"time"
)

// Event is a report of an event somewhere in the cluster. It generally denotes some state change in the system.
type V1beta1Event struct {

	// What action was taken/failed regarding to the regarding object.
	Action string `json:"action,omitempty"`

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources
	ApiVersion string `json:"apiVersion,omitempty"`

	// Deprecated field assuring backward compatibility with core.v1 Event type
	DeprecatedCount int32 `json:"deprecatedCount,omitempty"`

	// Deprecated field assuring backward compatibility with core.v1 Event type
	DeprecatedFirstTimestamp time.Time `json:"deprecatedFirstTimestamp,omitempty"`

	// Deprecated field assuring backward compatibility with core.v1 Event type
	DeprecatedLastTimestamp time.Time `json:"deprecatedLastTimestamp,omitempty"`

	// Deprecated field assuring backward compatibility with core.v1 Event type
	DeprecatedSource *V1EventSource `json:"deprecatedSource,omitempty"`

	// Required. Time when this Event was first observed.
	EventTime time.Time `json:"eventTime"`

	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	Kind string `json:"kind,omitempty"`

	Metadata *V1ObjectMeta `json:"metadata,omitempty"`

	// Optional. A human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB.
	Note string `json:"note,omitempty"`

	// Why the action was taken.
	Reason string `json:"reason,omitempty"`

	// The object this Event is about. In most cases it's an Object reporting controller implements. E.g. ReplicaSetController implements ReplicaSets and this event is emitted because it acts on some changes in a ReplicaSet object.
	Regarding *V1ObjectReference `json:"regarding,omitempty"`

	// Optional secondary object for more complex actions. E.g. when regarding object triggers a creation or deletion of related object.
	Related *V1ObjectReference `json:"related,omitempty"`

	// Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
	ReportingController string `json:"reportingController,omitempty"`

	// ID of the controller instance, e.g. `kubelet-xyzf`.
	ReportingInstance string `json:"reportingInstance,omitempty"`

	// Data about the Event series this event represents or nil if it's a singleton Event.
	Series *V1beta1EventSeries `json:"series,omitempty"`

	// Type of this event (Normal, Warning), new types could be added in the future.
	Type_ string `json:"type,omitempty"`
}
