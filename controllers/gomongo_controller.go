/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"context"

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ocp4eev1 "gitlab.com/ofekgit/go-operator-example/api/v1"
)

// GoMongoReconciler reconciles a GoMongo object
type GoMongoReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type getDeploymentFunc func(*GoMongoReconciler, *App, *ocp4eev1.GoMongo) *appsv1.Deployment

type App struct {
	Name          string
	Replicas      int32
	GetDeployment getDeploymentFunc
	Service       *Service
	Route         *Route
	Labels        map[string]string
}

type Service struct {
	Name       string
	Port       int32
	TargetPort int
	Labels     map[string]string
}

type Route struct {
	Name        string
	ServiceName string
}

// +kubebuilder:rbac:groups=ocp4ee.ocp4.example,resources=gomongoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ocp4ee.ocp4.example,resources=gomongoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ocp4ee.ocp4.example,resources=gomongoes/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments;services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=secrets;pods;pods/exec;pods/log;deployments;services,verbs=get;list;watch;create;update;patch;delete

func (r *GoMongoReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("gomongo", req.NamespacedName)

	// New operator instance
	gomongo := &ocp4eev1.GoMongo{}

	// Try to get existing operator instance
	err := r.Get(ctx, req.NamespacedName, gomongo)
	if err != nil {
		if errors.IsNotFound(err) {
			// Return and don't requeue
			log.Info("Gomongo instance not found. Ignoring since deletion of GoMongo instance means to stop operator.")
			return ctrl.Result{}, nil
		}
		// Any other error, requeue the request.
		log.Error(err, "Failed to get GoMongo instance")
		return ctrl.Result{}, err
	}

	apps := []App{
		{Name: "db",
			Replicas:      gomongo.Spec.DbSize,
			GetDeployment: getMongoDbDeployment,
			Labels:        labelsForMongoDb(),
			Service: &Service{
				Name:       "db-service",
				Port:       27017,
				TargetPort: 27017,
				Labels:     labelsForMongoDb(),
			},
			Route: nil,
		},
		{Name: "web",
			Replicas:      gomongo.Spec.WebSize,
			GetDeployment: getMongoWebDeployment,
			Labels:        labelsForMongoWeb(),
			Service: &Service{
				Name:       "mongoweb-service",
				Port:       8080,
				TargetPort: 8081,
				Labels:     labelsForMongoWeb(),
			},
			Route: &Route{
				Name:        "mongoweb-router",
				ServiceName: "mongoweb-service",
			},
		},
	}

	for _, app := range apps {
		dep := &appsv1.Deployment{}
		shouldReturn, result, err := r.ensureDeploymentExists(ctx, dep, &app, gomongo)
		if shouldReturn {
			return result, err
		}

		if shouldReturn, result, err := r.ensureDeploymentSize(ctx, dep, app.Replicas); shouldReturn {
			return result, err
		}

		if app.Service != nil {
			if shouldReturn, result, err := r.ensureServiceExists(ctx, &corev1.Service{}, &app, gomongo); shouldReturn {
				return result, err
			}
			if app.Route != nil {
				if shouldReturn, result, err := r.ensureRouteExists(ctx, &routev1.Route{}, &app, gomongo); shouldReturn {
					return result, err
				}
			}
		}

		podList := &corev1.PodList{}

		listOpts := []client.ListOption{
			client.InNamespace(gomongo.Namespace),
			client.MatchingLabels(app.Labels),
		}

		if err = r.List(ctx, podList, listOpts...); err != nil {
			log.Error(err, "Failed to list pods", "GoMongo.Namespace", gomongo.Namespace, "GoMongo.Name", gomongo.Name)
			return ctrl.Result{}, err
		}

		podNames := getPodNames(podList.Items)

		// Update status.Nodes if needed
		if !reflect.DeepEqual(podNames, gomongo.Status.Nodes) {
			gomongo.Status.Nodes = podNames
			err := r.Status().Update(ctx, gomongo)
			if err != nil {
				log.Error(err, "Failed to update GoMongo status")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *GoMongoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ocp4eev1.GoMongo{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&routev1.Route{}).
		Complete(r)
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
