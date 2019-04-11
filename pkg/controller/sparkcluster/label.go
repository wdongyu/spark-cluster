package sparkcluster

import (
	"fmt"
	sparkv1alpha1 "spark-cluster/pkg/apis/spark/v1alpha1"
)

func masterName(instance *sparkv1alpha1.SparkCluster) string {
	return instance.Spec.ClusterPrefix + "-" + Master
}

func masterLabel(instance *sparkv1alpha1.SparkCluster) map[string]string {
	return map[string]string{"app": masterName(instance)}
}

func masterPvc(instance *sparkv1alpha1.SparkCluster) string {
	return instance.Spec.ClusterPrefix + "-" + MasterPvc
}

func slaveName(instance *sparkv1alpha1.SparkCluster, index int) string {
	return instance.Spec.ClusterPrefix + "-" + Slave + "-" + fmt.Sprintf("%d", index)
}

func slaveLabel(instance *sparkv1alpha1.SparkCluster, index int) map[string]string {
	return map[string]string{"app": slaveName(instance, index)}
}

func slavePvc(instance *sparkv1alpha1.SparkCluster, index int) string {
	return instance.Spec.ClusterPrefix + "-" + SlavePvc + "-" + fmt.Sprintf("%d", index)
}
