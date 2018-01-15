package main

import (
	"flag"

	"context"
	log "github.com/sirupsen/logrus"
	"github.com/zmalik/k8s-publisher/pkg/controller"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
)

func main() {
	cfg, err := getConfig(*kubeconfig)

	if err != nil {
		log.Infoln("Error initializing the configuration : ", err)
		return
	}

	controller := controller.NewController(cfg)
	controller.Run(context.Background())

}

func init() {
	flag.Parse()
}

func getConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
