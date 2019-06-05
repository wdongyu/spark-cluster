package log

import (
	"io"

	"k8s.io/api/core/v1"
)

type LogDriver interface {
	// Name returns the name of log driver.
	Name() string

	// GetLog returns all raw log from a specified pod.
	GetLog(namespace string, pod v1.Pod) ([]string, error)

	// AggregateLogs aggregate all raw logs from given pod list.
	AggregateLogs(namespace string, pods []v1.Pod) (map[string][]string, error)

	// GetLogStream returns a read closer, which can be read to get the fully output of pod log.
	GetLogStream(namespace string, pod v1.Pod) (io.ReadCloser, error)

	// AggregateLogsStream returns a map of read closer, each connect with the pod and is used to read log
	// from container.
	AggregateLogStreams(namespace string, pods []v1.Pod) (map[string]io.ReadCloser, error)
}
