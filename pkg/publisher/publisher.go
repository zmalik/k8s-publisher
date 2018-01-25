package publisher

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/zmalik/k8s-publisher/pkg/publisher/slack"
	. "github.com/zmalik/k8s-publisher/pkg/publisher/types"
	"k8s.io/api/core/v1"
	"sync"
	"github.com/zmalik/k8s-publisher/pkg/publisher/http"
)

type Publisher struct {
}

type Notifiers []Notifier

type Notifier struct {
	Type  string `json:"type,omitempty" protobuf:"bytes,1,opt,name=type"`
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
}

func (p *Publisher) Process(key string, item interface{}) {
	if item == nil {
		// Pod was deleted
	} else {
		pod, ok := item.(*v1.Pod)
		if !ok {
			logrus.Errorf("Wrong object while adding the pod : %v", item)
		} else {
			// check the annotations to get the notify channels
			if channels, exists := p.getNotifyChannels(pod.Annotations); exists {
				// process rest of cases, based on pod phase
				if len(pod.Status.ContainerStatuses) > 0 {
					p.processPodStatusesInParallel(p.getPodStatus(key, pod), channels)
				}
			}
		}
	}
}

func (p *Publisher) getNotifyChannels(annotations map[string]string) (Notifiers, bool) {
	if notifiers, exists := annotations["notify-channels"]; exists {
		ntfiers := make(Notifiers, 0)
		err := json.Unmarshal([]byte(notifiers), &ntfiers)
		if err != nil {
			logrus.Errorf("Cannot unmarshall the annotation: %s, %v", notifiers, err)
		}
		return ntfiers, true
	}
	return nil, false

}

func (p *Publisher) processPodStatusesInParallel(podStatus *PodStatus, notifiersList Notifiers) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(notifiersList))

	defer waitGroup.Wait()

	for _, notifier := range notifiersList {
		go func() {
			defer waitGroup.Done()
			p.notify(notifier.Type, notifier.Value, podStatus)
		}()
	}
}

func (p *Publisher) getPodStatus(key string, pod *v1.Pod) *PodStatus {
	podStatus := &PodStatus{
		PodName: key,
		Phase:   string(pod.Status.Phase),
	}

	if pod.DeletionTimestamp != nil {
		podStatus.Phase = ToBeDeleted
		podStatus.DeletedAt = pod.DeletionTimestamp.Time
	}

	for _, status := range pod.Status.ContainerStatuses {
		if status.LastTerminationState.Terminated != nil {
			containerState := ContainerState{
				Image:    status.Image,
				State:    Terminated,
				ExitCode: status.LastTerminationState.Terminated.ExitCode,
				Reason:   status.LastTerminationState.Terminated.Reason,
				Message:  status.LastTerminationState.Terminated.Message,
			}
			podStatus.ContainersStates = append(podStatus.ContainersStates, containerState)
		}
		if status.LastTerminationState.Waiting != nil {
			containerState := ContainerState{
				Image:   status.Image,
				State:   Waiting,
				Reason:  status.LastTerminationState.Waiting.Reason,
				Message: status.LastTerminationState.Waiting.Message,
			}
			podStatus.ContainersStates = append(podStatus.ContainersStates, containerState)
		}

		if status.LastTerminationState.Running != nil {
			containerState := ContainerState{
				Image:     status.Image,
				State:     Running,
				StartedAt: status.LastTerminationState.Running.StartedAt.Time,
			}
			podStatus.ContainersStates = append(podStatus.ContainersStates, containerState)
		}
	}
	return podStatus
}

func (p *Publisher) notify(kind, value string, podStatus *PodStatus) {
	switch kind {
	// TODO add more plugins like stdout/email
	case "slack":
		slack.Publish(value, podStatus)
	case "http":
		http.Publish(value, podStatus)
	default:
		logrus.Warnf("Using a notifier type that isn't supported: %s with value %s", kind, value)
	}
}
