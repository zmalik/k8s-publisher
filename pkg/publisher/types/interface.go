package types

import "time"

const (
	Terminated  = "terminated"
	Running     = "running"
	Waiting     = "waiting"
	ToBeDeleted = "marked to be deleted"
)

type PodStatus struct {
	//pod Name with the format <namespace>/<name> unless namespace is empty
	PodName string `json:"podName,omitempty"`

	Phase string `json:"phase,omitempty"`

	ContainersStates []ContainerState `json:"containerStates,omitempty"`

	// Deletion timestamp, will be only populated it its phase is marked to be deleted
	DeletedAt time.Time `json:"deletedAt,omitempty"`
}

type ContainerState struct {
	// Can be running, terminated or waiting
	State string `json:"state,omitempty"`

	// Exit code, will be only populated it its terminated
	ExitCode int32 `json:"exitCode,omitempty"`

	// Container Image
	Image string `json:"image,omitempty"`

	// Reason of the state. Running state has no reason
	Reason string `json:"reason,omitempty"`

	// Message with the state. Running state has no message
	Message string `json:"message,omitempty"`

	// Container when it started. Only present with Running state
	StartedAt time.Time `json:"startedAt,omitempty"`
}
