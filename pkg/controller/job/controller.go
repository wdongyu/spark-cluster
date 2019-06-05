package job

import (
	sparkcluster "spark-cluster/pkg/apis/spark/v1alpha1"
	"spark-cluster/pkg/log"
	"spark-cluster/pkg/log/native"
	"spark-cluster/pkg/util/k8sutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type JobController struct {
	client     client.Client
	kubeClient kubernetes.Interface
	logDriver  log.LogDriver
}

func New() (*JobController, error) {
	clientset, err := setupClient()
	if err != nil {
		return nil, err
	}

	kubeClient, err := k8sutil.NewKubeClient()
	if err != nil {
		return nil, err
	}

	logDriver := native.New(kubeClient)
	if err != nil {
		return nil, err
	}

	return &JobController{
		client:     clientset,
		kubeClient: kubeClient,
		logDriver:  logDriver,
	}, nil
}

func setupClient() (client.Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	sparkcluster.AddToScheme(scheme)

	clientset, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (j *JobController) Execute() error {

	j.GetLogs()

	return nil
}
