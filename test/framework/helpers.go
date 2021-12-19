package framework

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"
	"time"
)

// GetFileDescriptor returns the file descriptor for the given relative file path.
func GetFileDescriptor(relativeFilePath string) (*os.File, error) {
	path, err := filepath.Abs(relativeFilePath)
	if err != nil {
		return nil, errors.Wrap(err,
			fmt.Sprintf("failed to generate absoulte path with %s", relativeFilePath))
	}
	return os.Open(path)
}

// WaitForPodReady checks whether expectedReplicas pods are running and ready.
func WaitForPodReady(client kubernetes.Interface, namespace string,
	timeout time.Duration, expectedReplicas int, opts metav1.ListOptions) error {

	return wait.Poll(time.Second, timeout, func() (done bool, err error) {
		podList, err := client.CoreV1().Pods(namespace).List(context.Background(), opts)
		if err != nil {
			return false, err
		}
		numRunningAndReadyPods := 0
		for _, pod := range podList.Items {
			ready, err := PodRunningAndReady(&pod)
			if err != nil {
				done = false
				return
			} else if ready == true {
				numRunningAndReadyPods++
			}
		}

		err = nil
		if numRunningAndReadyPods == expectedReplicas {
			done = true
			return
		}
		done = false
		return
	})
}
