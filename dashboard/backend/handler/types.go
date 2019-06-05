package handler

import (
	"spark-cluster/pkg/apis/spark/v1alpha1"
)

const (
	defaultNamespace       = "default"
	Token                  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0Ijoid2Rvbmd5dSIsInR5cGUiOiJ1c2VyIn0.315nHfICGYfYh71Nyv7xjIHLskJju1ZIm6ojWGuDvig"
	Host                   = "http://114.212.189.141:32015"
	Username               = "wdongyu"
	Email                  = "wdongyu@outlook.com"
	LogDriverElasticSearch = "elasticSearch"
	LogDriverNative        = "native"
)

type SparkClusterList struct {
	SparkClusters []v1alpha1.SparkCluster `json:"sparkclusters"`
}
