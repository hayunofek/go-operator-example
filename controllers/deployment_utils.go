package controllers

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"context"

	"github.com/prometheus/common/log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ocp4eev1 "gitlab.com/ofekgit/go-operator-example/api/v1"
)

func (r *GoMongoReconciler) ensureDeploymentExists(ctx context.Context, dep *appsv1.Deployment,
	app *App, operatorInstance *ocp4eev1.GoMongo) (bool, ctrl.Result, error) {

	err := r.Get(ctx, types.NamespacedName{Name: app.Name, Namespace: operatorInstance.Namespace}, dep)

	// If simply not found
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := app.GetDeployment(r, app, operatorInstance)
		controllerutil.SetControllerReference(operatorInstance, dep, r.Scheme)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return true, ctrl.Result{}, err
		}
		// Deployment created successfully - rerun Reconcile
		return true, ctrl.Result{Requeue: true}, nil

	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return true, ctrl.Result{}, err
	}

	return false, ctrl.Result{}, nil
}

func (r *GoMongoReconciler) ensureDeploymentSize(ctx context.Context, dep *appsv1.Deployment, size int32) (bool, ctrl.Result, error) {
	if *dep.Spec.Replicas != size {
		*dep.Spec.Replicas = size
		err := r.Update(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return true, ctrl.Result{}, err
		}

		return true, ctrl.Result{Requeue: true}, nil
	}

	return false, ctrl.Result{}, nil
}
