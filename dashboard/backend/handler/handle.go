package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"spark-cluster/pkg/apis"
	"spark-cluster/pkg/controller/job"
	"spark-cluster/pkg/log"
	"spark-cluster/pkg/log/native"
	"spark-cluster/pkg/util/k8sutil"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type APIHandler struct {
	frontDir string

	resourcesNamespace string

	kubeConfig    *rest.Config
	client        client.Client
	kubeClient    kubernetes.Interface
	logDriver     log.LogDriver
	jobController *job.JobController
}

func NewAPIHandler(frontDir string) (*APIHandler, error) {
	kubeConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// setup client set
	clientset, err := setupClient(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup kubernetes client: %v", err)
	}

	// setup kubernetes rest client
	kubeClient, err := k8sutil.NewKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to setup kubernetes client: %v", err)
	}

	jobController, err := job.New()
	if err != nil {
		return nil, err
	}

	apiHandler := &APIHandler{
		frontDir:      frontDir,
		client:        clientset,
		kubeClient:    kubeClient,
		kubeConfig:    kubeConfig,
		jobController: jobController,
	}

	// Set resources namespace
	apiHandler.resourcesNamespace = os.Getenv("RESOURCES_NAMESPACE")
	if len(apiHandler.resourcesNamespace) == 0 {
		apiHandler.resourcesNamespace = defaultNamespace
	}

	// Setup log driver
	logDriver := os.Getenv("LOG_DRIVER")
	logrus.Infof("Setting up log driver: %v", logDriver)
	if logDriver == LogDriverElasticSearch {
		return nil, errors.Wrap(err, "failed to setup log driver")
	} else if logDriver == LogDriverNative {
		apiHandler.logDriver = native.New(kubeClient)
	}

	return apiHandler, nil
}

func setupClient(config *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	for _, addToSchemeFunc := range []func(s *runtime.Scheme) error{
		apis.AddToScheme,
		v1.AddToScheme,
	} {
		if err := addToSchemeFunc(scheme); err != nil {
			return nil, err
		}
	}

	clientset, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

type Message struct {
	Message string `json:"message"`
}

func responseJSON(body interface{}, w http.ResponseWriter, statusCode int) {
	jsonResponse, err := json.Marshal(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
