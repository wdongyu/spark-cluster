package sparkcluster

import (
	"fmt"
	sparkv1alpha1 "spark-cluster/pkg/apis/spark/v1alpha1"
	"spark-cluster/pkg/controller/internal"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func SelectorForUser(user string) labels.Selector {
	selector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			internal.LabelUserKey: user,
		},
	}

	labelSelector, _ := metav1.LabelSelectorAsSelector(selector)

	return labelSelector
}

func masterName(instance *sparkv1alpha1.SparkCluster) string {
	return instance.Spec.ClusterPrefix + "-" + Master
}

func GetMasterLabel(instance *sparkv1alpha1.SparkCluster) map[string]string {
	return map[string]string{"app": masterName(instance)}
}

func masterPvc(instance *sparkv1alpha1.SparkCluster) string {
	return instance.Spec.ClusterPrefix + "-" + MasterPvc
}

func slaveName(instance *sparkv1alpha1.SparkCluster, index int) string {
	return instance.Spec.ClusterPrefix + "-" + Slave + "-" + fmt.Sprintf("%d", index)
}

func GetSlaveLabel(instance *sparkv1alpha1.SparkCluster, index int) map[string]string {
	return map[string]string{"app": slaveName(instance, index)}
}

func slavePvc(instance *sparkv1alpha1.SparkCluster, index int) string {
	return instance.Spec.ClusterPrefix + "-" + SlavePvc + "-" + fmt.Sprintf("%d", index)
}
