package native

import (
	"io"

	"strings"

	"spark-cluster/pkg/log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	nativeLogDriverName = "native"
)

type NativeLogDriver struct {
	client kubernetes.Interface
}

func New(client kubernetes.Interface) *NativeLogDriver {
	return &NativeLogDriver{
		client: client,
	}
}

func (d *NativeLogDriver) Name() string {
	return nativeLogDriverName
}

func (d *NativeLogDriver) GetLog(namespace string, pod v1.Pod) ([]string, error) {
	// check if the pod is running
	if pod.Status.Phase == v1.PodPending {
		return nil, log.ErrPodPending
	}

	raw, err := d.client.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{}).Do().Raw()
	if err != nil {
		return nil, nil
	}

	logs := strings.Split(string(raw), "\n")
	return logs, nil
}

func (d *NativeLogDriver) AggregateLogs(namespace string, pods []v1.Pod) (map[string][]string, error) {
	logs := make(map[string][]string)

	// check if all pods is not pending
	for _, pod := range pods {
		if pod.Status.Phase == v1.PodPending {
			return nil, log.ErrPodPending
		}

		l, err := d.GetLog(namespace, pod)
		if err != nil {
			return nil, err
		}
		logs[pod.Name] = l
	}

	return logs, nil
}

func (d *NativeLogDriver) GetLogStream(namespace string, pod v1.Pod) (io.ReadCloser, error) {
	logOpts := &v1.PodLogOptions{
		Follow: true,
	}

	// check if the pod is running
	if pod.Status.Phase == v1.PodPending {
		return nil, log.ErrPodPending
	}

	return d.client.CoreV1().Pods(namespace).GetLogs(pod.Name, logOpts).Stream()
}

func (d *NativeLogDriver) AggregateLogStreams(namespace string, pods []v1.Pod) (map[string]io.ReadCloser, error) {
	logOpts := &v1.PodLogOptions{
		Follow: true,
	}
	readers := make(map[string]io.ReadCloser)

	// check if all pods is running
	for _, pod := range pods {
		if pod.Status.Phase == v1.PodPending {
			return nil, log.ErrPodPending
		}

		readCloser, err := d.client.CoreV1().Pods(namespace).GetLogs(pod.Name, logOpts).Stream()
		if err != nil {
			return nil, err
		}

		readers[pod.Name] = readCloser
	}

	return readers, nil
}
