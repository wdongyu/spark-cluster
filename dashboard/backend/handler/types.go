package handler

import (
	"spark-cluster/pkg/apis/spark/v1alpha1"
)

const (
	Namespace = "default"
)

type SparkClusterList struct {
	SparkClusters []v1alpha1.SparkCluster `json:"sparkclusters"`
}
