/*
Copyright 2023.

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
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	podsmanagerv1 "pod-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PodsManageReconciler reconciles a PodsManage object
type PodsManageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=podsmanager.my.domain,resources=podsmanages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=podsmanager.my.domain,resources=podsmanages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=podsmanager.my.domain,resources=podsmanages/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodsManage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PodsManageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	podsManege := &podsmanagerv1.PodsManage{}
	if err := r.Get(ctx, req.NamespacedName, podsManege); err != nil {
		fmt.Println("the podsmanage was deleted!")
		return ctrl.Result{}, nil
	}

	deps := &appsv1.DeploymentList{}
	opts := []client.ListOption{
		client.InNamespace(req.NamespacedName.Namespace),
		client.MatchingLabels{"pods-manage": req.Name},
	}

	if err := r.List(ctx, deps, opts...); err != nil {
		fmt.Println(err)
		return ctrl.Result{}, err
	}

	for _, dep := range deps.Items {
		tempDep := dep
		tempContainers := make([]corev1.Container, 0)
		for _, container := range dep.Spec.Template.Spec.Containers {
			if container.Name == podsManege.Spec.ImageName {
				container.Image = podsManege.Spec.Image
				tempContainers = append(tempContainers, container)
			}
		}

		tempDep.Spec.Template.Spec.Containers = tempContainers

		if err := r.Update(ctx, &tempDep); err != nil {
			fmt.Println(err)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodsManageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podsmanagerv1.PodsManage{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
