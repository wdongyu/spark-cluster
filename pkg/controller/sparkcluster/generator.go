package sparkcluster

import (
	"fmt"
	sparkv1alpha1 "spark-cluster/pkg/apis/spark/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ReconcileSparkCluster) getMasterPod(instance *sparkv1alpha1.SparkCluster) *corev1.Pod {
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume

	if instance.Spec.NFS {
		nfs := corev1.NFSVolumeSource{Server: ShareServer, Path: "/hadoop/share-dir"}
		volumeMounts = append(volumeMounts, corev1.VolumeMount{Name: "share-dir", MountPath: nfs.Path})
		volumes = append(volumes, corev1.Volume{Name: "share-dir", VolumeSource: corev1.VolumeSource{NFS: &nfs}})
	}

	if instance.Spec.PvcEnable {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{Name: "dfs", MountPath: "/root/hdfs/namenode"})
		pvc := corev1.PersistentVolumeClaimVolumeSource{ClaimName: masterPvc(instance)}
		volumes = append(volumes, corev1.Volume{Name: "dfs", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &pvc}})
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      masterName(instance),
			Namespace: instance.Namespace,
			Labels:    GetMasterLabel(instance),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            Master,
					Image:           MasterImage,
					Command:         []string{"bash", "-c", "/etc/master-bootstrap.sh " + fmt.Sprintf("%d", instance.Spec.SlaveNum)},
					ImagePullPolicy: "IfNotPresent",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8020,
						},
						{
							ContainerPort: 50070,
						},
						{
							ContainerPort: 50470,
						},
					},
					Env:          []corev1.EnvVar{{Name: "PREFIX", Value: instance.Spec.ClusterPrefix}},
					Resources:    instance.Spec.Resources,
					VolumeMounts: volumeMounts,
				},
			},
			Volumes: volumes,
		},
	}
}

func (r *ReconcileSparkCluster) getSlavePod(instance *sparkv1alpha1.SparkCluster, index int) *corev1.Pod {
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume
	if instance.Spec.PvcEnable {
		volumeMounts = []corev1.VolumeMount{{Name: "dfs", MountPath: "/root/hdfs/datanode"}}
		pvc := corev1.PersistentVolumeClaimVolumeSource{ClaimName: slavePvc(instance, index)}
		volumes = []corev1.Volume{{Name: "dfs", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &pvc}}}
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      slaveName(instance, index),
			Namespace: instance.Namespace,
			Labels:    GetSlaveLabel(instance, index),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            Slave,
					Image:           SlaveImage,
					ImagePullPolicy: "IfNotPresent",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 50010,
						},
						{
							ContainerPort: 50020,
						},
						{
							ContainerPort: 50075,
						},
						{
							ContainerPort: 50475,
						},
					},
					Env:          []corev1.EnvVar{{Name: "PREFIX", Value: instance.Spec.ClusterPrefix}},
					Resources:    instance.Spec.Resources,
					VolumeMounts: volumeMounts,
				},
			},
			Volumes: volumes,
		},
	}
}

func (r *ReconcileSparkCluster) getMasterService(instance *sparkv1alpha1.SparkCluster) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      masterName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name: "rpc",
					Port: 8020,
				},
				{
					Name: "p1",
					Port: 50070,
				},
				{
					Name: "p2",
					Port: 50470,
				},
			},
			Selector: GetMasterLabel(instance),
		},
	}
}

func (r *ReconcileSparkCluster) getUIService(instance *sparkv1alpha1.SparkCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name: "ssh",
			Port: 22,
		},
		{
			Name: "hdfs-client",
			Port: 9000,
		},
		{
			Name: "resource-manager",
			Port: 8088,
		},
		{
			Name: "name-node",
			Port: 50070,
		},
		{
			Name: "spark",
			Port: 8080,
		},
		{
			Name: "spark-shell",
			Port: 4040,
		}}
	ports = append(ports, instance.Spec.Ports...)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.ClusterPrefix + "-ui-service",
			Namespace: instance.Namespace,
			Labels:    GetMasterLabel(instance),
		},
		Spec: corev1.ServiceSpec{
			Type:     "NodePort",
			Ports:    ports,
			Selector: GetMasterLabel(instance),
		},
	}
}

func (r *ReconcileSparkCluster) getSlaveService(instance *sparkv1alpha1.SparkCluster, index int) *corev1.Service {
	serviceName := slaveName(instance, index)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name: "rpc",
					Port: 8020,
				},
				{
					Name: "p1",
					Port: 50070,
				},
				{
					Name: "p2",
					Port: 50470,
				},
			},
			Selector: map[string]string{"app": serviceName},
		},
	}
}

func (r *ReconcileSparkCluster) getMasterPvc(instance *sparkv1alpha1.SparkCluster) *corev1.PersistentVolumeClaim {
	storageClassName := StorageClassName
	accessModes := make([]corev1.PersistentVolumeAccessMode, 1)
	accessModes[0] = corev1.ReadWriteOnce
	resourceList := corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("20Gi")}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MasterPvc,
			Namespace: instance.Namespace,
			Labels:    GetMasterLabel(instance),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      accessModes,
			Resources:        corev1.ResourceRequirements{Requests: resourceList},
		},
	}
}

func (r *ReconcileSparkCluster) getSlavePvc(instance *sparkv1alpha1.SparkCluster, index int) *corev1.PersistentVolumeClaim {
	storageClassName := StorageClassName
	accessModes := make([]corev1.PersistentVolumeAccessMode, 1)
	accessModes[0] = corev1.ReadWriteOnce
	q := resource.MustParse("20Gi")
	resourceList := corev1.ResourceList{corev1.ResourceStorage: q}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      slavePvc(instance, index),
			Namespace: instance.Namespace,
			Labels:    GetSlaveLabel(instance, index),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      accessModes,
			Resources:        corev1.ResourceRequirements{Requests: resourceList},
		},
	}
}
