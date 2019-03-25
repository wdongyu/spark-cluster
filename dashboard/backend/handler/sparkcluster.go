package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"io/ioutil"

	sparkclusterv1alpha1 "spark-cluster/pkg/apis/spark/v1alpha1"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListSparkCluster(w http.ResponseWriter, r *http.Request) {
	sc := new(sparkclusterv1alpha1.SparkClusterList)

	opts := &client.ListOptions{}
	opts.SetLabelSelector(fmt.Sprintf("app=%s", "hadoop-spark-cluster"))
	opts.InNamespace(Namespace)
	err := handler.client.List(context.TODO(), opts, sc)

	if err != nil {
		log.Warningf("failed to list spark cluster: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(SparkClusterList{SparkClusters: sc.Items}, w, http.StatusOK)
	}
}

func (handler *APIHandler) CreateSparkCluster(w http.ResponseWriter, r *http.Request) {
	sc := new(sparkclusterv1alpha1.SparkCluster)
	// user := r.Header.Get("User")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &sc); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}

	if len(sc.Namespace) == 0 {
		sc.Namespace = Namespace
	}
	// workspace.AddUserLabel(ws, user)

	err = handler.client.Create(context.TODO(), sc)
	if err != nil {
		log.Warningf("Failed to create spark cluster %v: %v", sc.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(sc, w, http.StatusCreated)
	}
}
