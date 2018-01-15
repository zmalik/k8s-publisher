package slack

import (
	"bytes"
	"fmt"
	slack_api "github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/zmalik/k8s-publisher/pkg/publisher/types"
	"os"
	"strings"
	"time"
)

var (
	client   *slack_api.Client
	userName string
)

func init() {
	userName = os.Getenv("SLACK_USERNAME")
	token := strings.TrimSuffix(os.Getenv("SLACK_API_TOKEN"), "\n")
	if len(token) == 0 {
		log.Panicf("Missing SLACK_API_TOKEN environment variable.")
	}
	client = slack_api.New(token)
}

func Publish(channel string, podStatus *types.PodStatus) {
	params := slack_api.PostMessageParameters{
		Username: userName,
	}
	channelID, timestamp, err := client.PostMessage(channel, formatMessage(podStatus), params)
	if err != nil {
		log.Errorf("Error sending messages to the %s : %s", channel, err.Error())
	}
	log.Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func formatMessage(podStatus *types.PodStatus) string {
	var buffer bytes.Buffer
	switch podStatus.Phase {
	case types.ToBeDeleted:
		buffer.WriteString(fmt.Sprintf("Pod: *%s %s* for %s \n",
			podStatus.PodName, podStatus.Phase, podStatus.DeletedAt.Format(time.UnixDate)))
	default:
		buffer.WriteString(fmt.Sprintf("Pod: *%s %s* \n",
			podStatus.PodName, podStatus.Phase))
	}

	for _, containerState := range podStatus.ContainersStates {
		switch containerState.State {
		case types.Running:
			buffer.WriteString(fmt.Sprintf("container *%s* started *running* at: %s \n",
				containerState.Image, containerState.StartedAt.Format(time.UnixDate)))
		case types.Waiting:
			buffer.WriteString(fmt.Sprintf("container *%s is waiting*. Reason: %s Message: %s\n",
				containerState.Image, containerState.Reason,
				containerState.Message))
		case types.Terminated:
			buffer.WriteString(fmt.Sprintf("container *%s terminated*. Code: %d Reason: %s Message: %s\n",
				containerState.Image, containerState.ExitCode,
				containerState.Reason, containerState.Message))
		}
	}
	return buffer.String()
}
