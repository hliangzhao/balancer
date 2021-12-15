package balancer

import (
	`context`
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	appv1 `k8s.io/api/apps/v1`
	corev1 `k8s.io/api/core/v1`
	`k8s.io/apimachinery/pkg/api/errors`
	metav1 `k8s.io/apimachinery/pkg/apis/meta/v1`
	`k8s.io/apimachinery/pkg/types`
	`sigs.k8s.io/controller-runtime/pkg/controller/controllerutil`
)

func (r *ReconcileBalancer) syncDeployment(balancer *balancerv1alpha1.Balancer) error {
	// 1 firstly, we sync configmap
	cm, err := r.syncConfigMap(balancer)
	if err != nil {
		return err
	}

	// 2 now we sync deployment
	dp, err := NewDeployment(balancer)
	if err != nil {
		return err
	}
	annotations := map[string]string{
		balancerv1alpha1.ConfigMapHashKey: ConfigMapHash(cm),
	}
	// TODO: the assignment should happen only when cm changes
	dp.Spec.Template.ObjectMeta.Annotations = annotations

	// set balancer as the owner of dp
	if err = controllerutil.SetOwnerReference(balancer, dp, r.scheme); err != nil {
		return err
	}

	found := &appv1.Deployment{}
	err = r.client.Get(context.Background(), types.NamespacedName{Namespace: balancer.Namespace, Name: balancer.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		// corresponding dp not found in the cluster, create it with the newest dp
		if err = r.client.Create(context.Background(), dp); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// corresponding dp found, update it with the newest dp
	found.Spec.Template = dp.Spec.Template
	if err = r.client.Update(context.Background(), found); err != nil {
		return  err
	}
	return nil
}

// TODO: What about multiple pods? Relations with the number of endpoints?

// NewDeployment creates a new deployment (with one nginx pod) for the Balancer.
func NewDeployment(balancer *balancerv1alpha1.Balancer) (*appv1.Deployment, error) {
	replicas := int32(1)
	labels := NewPodLabels(balancer)
	nginxContainer := corev1.Container{
		Name:  "nginx",
		Image: "nginx:1.15.9",
		Ports: []corev1.ContainerPort{{ContainerPort: 80}},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      ConfigMapName(balancer),
				MountPath: "/etc/nginx",
				ReadOnly:  true,
			},
		},
	}
	nginxVolume := corev1.Volume{
		Name: ConfigMapName(balancer),
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: ConfigMapName(balancer),
			},
		}},
	}

	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeploymentName(balancer),
			Namespace: balancer.Namespace,
			Labels:    labels,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      DeploymentName(balancer),
					Namespace: balancer.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{nginxContainer},
					Volumes:    []corev1.Volume{nginxVolume},
				},
			},
		},
		Status: appv1.DeploymentStatus{},
	}, nil
}

func DeploymentName(balancer *balancerv1alpha1.Balancer) string {
	return balancer.Name + "proxy"
}
