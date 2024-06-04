/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logger "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"

	externalresourcedevv1alpha1 "github.com/Y0sh1dk/kubebuilder-external-resource-demo/api/v1alpha1"
	todoClient "github.com/Y0sh1dk/kubebuilder-external-resource-demo/internal/clients/todo"
	"github.com/Y0sh1dk/kubebuilder-external-resource-demo/internal/todo"
	"github.com/k0kubun/pp/v3"
	"github.com/samber/lo"
)

const (
	deleteFinalizer = "finalizer.external-resource.dev"
)

// TodoReconciler reconciles a Todo object
type TodoReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	TodoClient *todoClient.Client
}

//+kubebuilder:rbac:groups=external-resource.dev.external-resource.dev,resources=todoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=external-resource.dev.external-resource.dev,resources=todoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=external-resource.dev.external-resource.dev,resources=todoes/finalizers,verbs=update

func (r *TodoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx)
	log.Info("Reconciling Todo", "Name", req.NamespacedName, "Namespace", req.Namespace)

	t := &externalresourcedevv1alpha1.Todo{}
	if err := r.Get(ctx, req.NamespacedName, t); err != nil {
		log.Error(err, "Failed to get Todo")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle deletion
	if t.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(t, deleteFinalizer) {
			controllerutil.AddFinalizer(t, deleteFinalizer)
			if err := r.Update(ctx, t); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(t, deleteFinalizer) {
			if err := r.TodoClient.DeleteTodo(t.Status.ID); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(t, deleteFinalizer)
			if err := r.Update(ctx, t); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	todoBackendObj := &todo.Todo{
		ID:    t.Status.ID,
		Title: t.Spec.Title,
	}

	err := r.createOrUpdate(ctx, todoBackendObj)
	if err != nil {
		log.Error(err, "Failed to create or update Todo")
		return ctrl.Result{}, err
	}

	pp.Println(todoBackendObj)

	t.Status.ID = todoBackendObj.ID
	if err := r.Status().Update(ctx, t); err != nil {
		log.Error(err, "Failed to update Todo status")

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TodoReconciler) SetupWithManager(mgr ctrl.Manager, client *todoClient.Client) error {
	externalEventChan := make(chan event.GenericEvent)

	mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		for {
			fmt.Println("Getting all ToDos from API")
			apiTodos, err := client.GetTodos()
			if err != nil {
				continue
			}

			fmt.Println("Getting all ToDos from Kube")
			kubeTodos := &externalresourcedevv1alpha1.TodoList{}
			if err := mgr.GetClient().List(ctx, kubeTodos); err != nil {
				continue
			}

			lo.ForEach(apiTodos, func(item todo.Todo, index int) {
				lo.ForEach(kubeTodos.Items, func(kubeTodo externalresourcedevv1alpha1.Todo, index int) {
					if kubeTodo.Status.ID == item.ID {
						return
					}

					fmt.Println("Adding Todo to workqueue", "name", kubeTodo.Name, "namespace", kubeTodo.Namespace)
					externalEventChan <- event.GenericEvent{
						Object: &kubeTodo,
					}
				})
			})

			time.Sleep(5 * time.Second)
		}
	}))

	return ctrl.NewControllerManagedBy(mgr).
		For(&externalresourcedevv1alpha1.Todo{}).
		WatchesRawSource(&source.Channel{Source: externalEventChan}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

func (r *TodoReconciler) createOrUpdate(ctx context.Context, t *todo.Todo) error {
	log := logger.FromContext(ctx)

	existing, err := r.TodoClient.GetTodo(t.ID)
	if err != nil { // Does not exist, create it
		log.Info("Creating Todo", "Title", t.Title)

		created, err := r.TodoClient.CreateTodo(*t)
		if err != nil {
			log.Info("Failed to create Todo", "Title", t.Title)

			return err
		}

		t.ID = created.ID

		return nil
	}

	// Does exist, update it
	log.Info("Updating Todo", "Title", t.Title)
	existing.Title = t.Title

	updated, err := r.TodoClient.UpdateTodo(*existing)
	if err != nil {
		log.Error(err, "Failed to update Todo")

		return err
	}

	t.ID = updated.ID
	t.Title = updated.Title

	return nil
}
