package client

import (
	"spark-cluster/pkg/util/k8sutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ClientManager struct {
	restCfg   *rest.Config
	ClientSet *kubernetes.Clientset
}

func (c *ClientManager) init() error {
	restCfg, err := k8sutil.GetClusterConfig()
	if err != nil {
		return err
	}
	c.restCfg = restCfg

	clientset, err := kubernetes.NewForConfig(c.restCfg)
	if err != nil {
		return err
	}
	c.ClientSet = clientset

	return nil
}

func NewClientManager() (ClientManager, error) {
	cm := ClientManager{}
	err := cm.init()

	if err != nil {
		return ClientManager{}, err
	}

	return cm, err
}
