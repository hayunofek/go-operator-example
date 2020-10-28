package controllers

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ocp4eev1 "gitlab.com/ofekgit/go-operator-example/api/v1"
)

func getMongoDbDeployment(r *GoMongoReconciler, app *App, m *ocp4eev1.GoMongo) *appsv1.Deployment {
	replicas := m.Spec.DbSize
	ls := labelsForMongoDb()
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "registry.redhat.io/rhscl/mongodb-34-rhel7",
						Name:  "mongodb",
						Env: []corev1.EnvVar{{
							Name:  "MONGODB_ADMIN_PASSWORD",
							Value: m.Spec.MongoDbAdminPass,
						}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 27017,
						}},
					}},
				},
			},
		},
	}
}

func labelsForMongoDb() map[string]string {
	return map[string]string{"app": "mongodb"}
}
