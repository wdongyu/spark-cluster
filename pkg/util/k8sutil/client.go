package k8sutil

import (
	"os"
	"path/filepath"

	"spark-cluster/pkg/util"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// RecommendedConfigPathEnvVar is a environment variable for path configuration
const RecommendedConfigPathEnvVar = "KUBECONFIG"

// GetClusterConfig obtain the config from the Kube configuration used by kubeconfig, or from k8s cluster.
func GetClusterConfig() (*rest.Config, error) {
	if len(os.Getenv(RecommendedConfigPathEnvVar)) > 0 {
		// use the current context in kubeconfig
		return clientcmd.BuildConfigFromFlags("", os.Getenv(RecommendedConfigPathEnvVar))
	}

	if home := homedir.HomeDir(); home != "" {
		// use the kubeconfig in $HOME/.kube/config
		kubeconfig := filepath.Join(home, ".kube", "config")
		if util.IsExist(kubeconfig) {
			return clientcmd.BuildConfigFromFlags("", kubeconfig)
		}
	}

	return rest.InClusterConfig()
}

// NewKubeClient returns new kubernetes client for cluster configuration
func NewKubeClient() (*kubernetes.Clientset, error) {
	cfg, err := GetClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}
