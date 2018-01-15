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
	PodName string

	Phase string

	ContainersStates []ContainerState

	// Deletion timestamp, will be only populated it its phase is marked to be deleted
	DeletedAt time.Time
}

type ContainerState struct {
	// Can be running, terminated or waiting
	State string

	// Exit code, will be only populated it its terminated
	ExitCode int32

	// Container Image
	Image string

	// Reason of the state. Running state has no reason
	Reason string

	// Message with the state. Running state has no message
	Message string

	// Container when it started. Only present with Running state
	StartedAt time.Time
}
