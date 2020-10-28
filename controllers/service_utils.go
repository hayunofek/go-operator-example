package controllers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	"context"

	"github.com/prometheus/common/log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ocp4eev1 "gitlab.com/ofekgit/go-operator-example/api/v1"
)

func (r *GoMongoReconciler) ensureServiceExists(ctx context.Context, service *corev1.Service,
	app *App, operatorInstance *ocp4eev1.GoMongo) (bool, ctrl.Result, error) {

	err := r.Get(ctx, types.NamespacedName{Name: app.Service.Name, Namespace: operatorInstance.Namespace}, service)

	// If simply not found
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Service.Name,
				Namespace: operatorInstance.Namespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: app.Labels,
				Ports: []corev1.ServicePort{
					{
						Port:       app.Service.Port,
						TargetPort: intstr.FromInt(app.Service.TargetPort),
					},
				},
			},
		}
		controllerutil.SetControllerReference(operatorInstance, service, r.Scheme)
		log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.Create(ctx, service)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return true, ctrl.Result{}, err
		}
		// Service created successfully - rerun Reconcile
		return true, ctrl.Result{Requeue: true}, nil

	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return true, ctrl.Result{}, err
	}

	return false, ctrl.Result{}, nil
}
