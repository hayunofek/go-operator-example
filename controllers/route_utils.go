package controllers

import (
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"context"

	"github.com/prometheus/common/log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ocp4eev1 "gitlab.com/ofekgit/go-operator-example/api/v1"
)

func (r *GoMongoReconciler) ensureRouteExists(ctx context.Context, route *routev1.Route,
	app *App, operatorInstance *ocp4eev1.GoMongo) (bool, ctrl.Result, error) {
	err := r.Get(ctx, types.NamespacedName{Name: app.Route.Name, Namespace: operatorInstance.Namespace}, route)

	// If simply not found
	if err != nil && errors.IsNotFound(err) {
		// Define a new route
		route := &routev1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Route.Name,
				Namespace: operatorInstance.Namespace,
			},
			Spec: routev1.RouteSpec{
				To: routev1.RouteTargetReference{
					Kind: "Service",
					Name: app.Route.ServiceName,
				},
			},
		}
		controllerutil.SetControllerReference(operatorInstance, route, r.Scheme)
		log.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		err = r.Create(ctx, route)
		if err != nil {
			log.Error(err, "Failed to create new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
			return true, ctrl.Result{}, err
		}
		// Service created successfully - rerun Reconcile
		return true, ctrl.Result{Requeue: true}, nil

	} else if err != nil {
		log.Error(err, "Failed to get Route")
		return true, ctrl.Result{}, err
	}

	return false, ctrl.Result{}, nil
}
