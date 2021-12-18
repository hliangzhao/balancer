package framework

import (
	`context`
	`fmt`
	`github.com/pkg/errors`
	corev1 `k8s.io/api/core/v1`
	metav1 `k8s.io/apimachinery/pkg/apis/meta/v1`
	`k8s.io/apimachinery/pkg/util/wait`
	`k8s.io/apimachinery/pkg/util/yaml`
	`time`
)

// CreatePod creates pod in given namespace.
func (f *Framework) CreatePod(namespace string, pod *corev1.Pod) error {
	pod.Namespace = namespace
	_, err := f.KubeClient.CoreV1().Pods(namespace).Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// WaitForPodReady detects pod is ready or not within timeout.
func (f *Framework) WaitForPodReady(pod *corev1.Pod, timeout time.Duration) error {
	var pollErr error
	err := wait.Poll(2*time.Second, timeout, func() (bool, error) {
		pod, err := f.KubeClient.CoreV1().Pods(pod.Namespace).Get(context.Background(), pod.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if pod.Status.Phase != corev1.PodRunning {
			return false, nil
		}
		return true, nil
	})
	return errors.Wrapf(pollErr, "waiting for Balancer %s/%s: %v", pod.Namespace, pod.Name, err)
}

func (f *Framework) CreatePodAndWaitUntilReady(namespace string, pod *corev1.Pod) error {
	if err := f.CreatePod(namespace, pod); err != nil {
		return err
	}
	if err := f.WaitForPodReady(pod, 30*time.Second); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create pod %s", pod.Name))
	}
	return nil
}

// GetPodRestartCount gets a map recording each container's restart count for the given pod.
func (f *Framework) GetPodRestartCount(namespace, name string) (map[string]int32, error) {
	pod, err := f.KubeClient.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	restarts := map[string]int32{}
	for _, status := range pod.Status.ContainerStatuses {
		restarts[status.Name] = status.RestartCount
	}
	return restarts, nil
}

// PodRunningAndReady returns true if pod is ready and its current status is running.
func PodRunningAndReady(pod *corev1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case corev1.PodFailed, corev1.PodSucceeded:
		return false, fmt.Errorf("pod failed or completed")
	case corev1.PodRunning:
		// a pod could experience many conditions, we need to judge the `PodReady` condition exists,
		// and it is the current condition exactly
		for _, cond := range pod.Status.Conditions {
			if cond.Type != corev1.PodReady {
				continue
			}
			return cond.Status == corev1.ConditionTrue, nil
		}
		return false, fmt.Errorf("pod ready condition not found")
	}
	return false, fmt.Errorf("unexpected error")
}

func ParsePodYaml(relativeFilePath string) (*corev1.Pod, error) {
	manifest, err := GetFileDescriptor(relativeFilePath)
	if err != nil {
		return nil, err
	}
	pod := &corev1.Pod{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(pod); err != nil {
		return nil, err
	}
	return pod, nil
}
