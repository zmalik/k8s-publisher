package http

import (
	"github.com/zmalik/k8s-publisher/pkg/publisher/types"
	"net/http"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"bytes"
)

func Publish(endpoint string, podStatus *types.PodStatus) {
	body, err := json.Marshal(podStatus)
	if err != nil{
		log.Errorf("Error marshaling the pod status: %s", err.Error())
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil{
		log.Errorf("Error posting the pod status: %s", err.Error())
	}
	if resp.StatusCode >= 300 {
		log.Errorf("Expected: {1,2}xx publishing the pod status, got: %d", resp.StatusCode)
	}
}
