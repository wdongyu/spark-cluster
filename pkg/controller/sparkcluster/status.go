package sparkcluster

import (
	"context"
	"fmt"
	"reflect"
	"spark-cluster/pkg/apis/spark/v1alpha1"
	sparkv1alpha1 "spark-cluster/pkg/apis/spark/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileSparkCluster) updateStatus(instance *sparkv1alpha1.SparkCluster) error {
	if instance.Status.CreateTime == nil {
		now := metav1.Now()
		instance.Status.CreateTime = &now
	}

	pods, err := PodsForLabels(instance, r.Client)
	if err != nil {
		log.Error(err, "failed to update pod status for SparkCluster %v", instance.Name)
		return err
	}

	// services, err := ServicesForLabels(instance, r.Client)
	uiService := &corev1.Service{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: UIService, Namespace: instance.Namespace}, uiService)
	if err != nil {
		log.Error(err, "failed to update service status for SparkCluster %v", instance.Name)
		return err
	}

	podStatuses := MappingPodsByPhase(pods)
	if podStatuses[corev1.PodRunning] >= instance.Spec.SlaveNum+1 {

		// if all pods are running, we set the SparkCluster status to running and update endpoint field
		instance.Status.Phase = sparkv1alpha1.SparkClusterPhaseRunning
		if err := r.updateEndpoints(instance, pods, *uiService); err != nil {
			log.Error(err, "failed to update pod endpoints")
		}

	} else if podStatuses[corev1.PodFailed] > 0 {
		instance.Status.Phase = sparkv1alpha1.SparkClusterPhaseFailed
	} else {
		instance.Status.Phase = sparkv1alpha1.SparkClusterPhasePending
	}

	return r.syncStatus(instance)
}

func (r *ReconcileSparkCluster) syncStatus(instance *v1alpha1.SparkCluster) error {
	old := &v1alpha1.SparkCluster{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}, old)
	if err != nil {
		log.Error(err, "failed to update status for SparkCluster %v", instance.Name)
		return err
	}

	if !reflect.DeepEqual(old.Status, instance.Status) {
		return r.Update(context.TODO(), instance)
	}

	return nil
}

func (r *ReconcileSparkCluster) updateEndpoints(instance *v1alpha1.SparkCluster, pods []corev1.Pod, uiService corev1.Service) error {
	if instance.Status.Endpoints == nil {
		instance.Status.Endpoints = make(map[string]string)
	}
	if instance.Status.ExposedPorts == nil {
		instance.Status.ExposedPorts = make([]corev1.ServicePort, len(uiService.Spec.Ports))
	}

	var podName string
	for i, pod := range pods {
		if i == 0 {
			podName = MasterName
		} else {
			podName = SlaveName + "-" + fmt.Sprintf("%d", i)
		}
		instance.Status.Endpoints[podName] = pod.Status.PodIP
	}

	instance.Status.ExposedPorts = uiService.Spec.Ports
	return nil
}

func PodsForLabels(instance *sparkv1alpha1.SparkCluster, c client.Client) ([]corev1.Pod, error) {
	pods := make([]corev1.Pod, instance.Spec.SlaveNum+1)

	masterPod := &corev1.Pod{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: MasterName, Namespace: instance.Namespace}, masterPod)
	if err != nil {
		return nil, err
	}
	pods[0] = *masterPod
	// pods = append(pods, *masterPod)

	for i := 1; i <= instance.Spec.SlaveNum; i++ {
		pod := &corev1.Pod{}
		err := c.Get(context.TODO(), types.NamespacedName{Name: SlaveName + "-" + fmt.Sprintf("%d", i), Namespace: instance.Namespace}, pod)
		if err != nil {
			return nil, err
		}
		// pods = append(pods, *pod)
		pods[i] = *pod
	}

	return pods, nil
}

func ServicesForLabels(instance *sparkv1alpha1.SparkCluster, c client.Client) ([]corev1.Service, error) {
	services := make([]corev1.Service, instance.Spec.SlaveNum+2)

	masterService := &corev1.Service{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: MasterName, Namespace: instance.Namespace}, masterService)
	if err != nil {
		return nil, err
	}
	services = append(services, *masterService)

	uiService := &corev1.Service{}
	err = c.Get(context.TODO(), types.NamespacedName{Name: UIService, Namespace: instance.Namespace}, uiService)
	if err != nil {
		return nil, err
	}
	services = append(services, *uiService)

	for i := 1; i <= instance.Spec.SlaveNum; i++ {
		service := &corev1.Service{}
		err := c.Get(context.TODO(), types.NamespacedName{Name: SlaveName + "-" + fmt.Sprintf("%d", i), Namespace: instance.Namespace}, service)
		if err != nil {
			return nil, err
		}
		services = append(services, *service)
	}

	return services, nil
}

func MappingPodsByPhase(pods []corev1.Pod) map[corev1.PodPhase]int {
	result := make(map[corev1.PodPhase]int)
	for _, pod := range pods {
		if len(pod.Status.Phase) == 0 {
			continue
		}
		if _, ok := result[pod.Status.Phase]; !ok {
			result[pod.Status.Phase] = 1
		} else {
			result[pod.Status.Phase]++
		}
	}
	return result
}
