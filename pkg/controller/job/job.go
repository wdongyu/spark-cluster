package job

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	maxRetry      = 5
	queryInterval = 5 * time.Second
)

// GetLogs gathers all logs from job pods and print them to STDOUT
func (j *JobController) GetLogs() {
	var pods *v1.PodList
	var wg sync.WaitGroup
	var err error
	var name = "spark-hive"

	fmt.Printf("[sparkctl] Setting up log driver: %v\n", j.logDriver.Name())
	fmt.Printf("[sparkctl] Getting Job %s\n", name)

	retry := 0
	podMap := make(map[string]v1.Pod)
	for {
		pods, err = j.kubeClient.CoreV1().Pods("kerong").List(metav1.ListOptions{
			LabelSelector: "driver-pod=spark-hive",
		})
		if err != nil {
			fmt.Println(err)
			fmt.Printf("[sparkctl] Failed to get pods for the given job %s\n", name)
			if retry == maxRetry {
				fmt.Printf("[sparkctl] Tried %d times but cannot get the pods\n", maxRetry)
				return
			}
			retry++
			continue
		}
		if checkPodStatus(pods) {
			fmt.Printf("[sparkctl] Waiting for creating pods for the given job\n")
			time.Sleep(queryInterval)
			continue
		}
		for _, pod := range pods.Items {
			if pod.Status.Phase != v1.PodPending {
				podMap[pod.Name] = pod
			}
		}
		if len(podMap) != 1 {
			fmt.Printf("[sparkctl] Waiting for running all pod (%d master)\n", 1)
			time.Sleep(queryInterval)
		} else {
			break
		}
	}

	fmt.Printf("[sparkctl] There are %d pods for the job.\n", len(pods.Items))

	var readerClosers map[string]io.ReadCloser
	var retryInterval = queryInterval
	for {
		readerClosers, err = j.logDriver.AggregateLogStreams("kerong", pods.Items)
		// readerClosers, err = j.logDriver.AggregateLogs("default", pods.Items)
		if err != nil {
			fmt.Printf("[sparkctl] Failed to get pod log: %v, retry in %v seconds\n", err, retryInterval.Seconds())
			if retry == maxRetry {
				fmt.Printf("[sparkctl] Tried %d times but cannot get the pods\n", maxRetry)
				return
			}
			retry++
			retryInterval *= 2
			time.Sleep(retryInterval)
			continue
		} else {
			break
		}
	}

	wg.Add(len(pods.Items))
	for podName, closer := range readerClosers {
		podRef := podMap[podName]

		go func(closer io.ReadCloser, pod v1.Pod) {
			defer closer.Close()

			reader := bufio.NewReader(closer)
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Printf("[sparkctl] Failed to read log: %v\n", err)
					break
				}

				fmt.Printf("[sparkctl] %s", string(line))
			}
			wg.Done()
		}(closer, podRef)
	}

	wg.Wait()
	fmt.Printf("[sparkctl] Finished\n")
	return
}

func checkPodStatus(pods *v1.PodList) bool {
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodSucceeded {
			return true
		}
	}

	return false
}
