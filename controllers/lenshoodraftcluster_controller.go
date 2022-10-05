/*
Copyright 2022.

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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	. "go-raft-operator/api/v1"
)

// LenshoodRaftClusterReconciler reconciles a LenshoodRaftCluster object
type LenshoodRaftClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lenshood.github.io,resources=lenshoodraftclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lenshood.github.io,resources=lenshoodraftclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lenshood.github.io,resources=lenshoodraftclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LenshoodRaftCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *LenshoodRaftClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
	_ = log.FromContext(ctx)

	cluster := &LenshoodRaftCluster{}
	err = r.Client.Get(ctx, req.NamespacedName, cluster)
	if err != nil {
		if errors.IsNotFound(err) {
			return res, deleteCluster(cluster.DeepCopy())
		}
		return
	}

	if reachedDesiredState(cluster) {
		return
	}

	if len(cluster.Status.Ids) == 0 {
		return res, createSeedToCluster(cluster.DeepCopy())
	}

	return ctrl.Result{}, joinToCluster(cluster.DeepCopy())
}

func reachedDesiredState(cluster *LenshoodRaftCluster) bool {
	return false
}

func createSeedToCluster(cluster *LenshoodRaftCluster) error {
	return nil
}

func joinToCluster(cluster *LenshoodRaftCluster) error {
	return nil
}

func deleteCluster(cluster *LenshoodRaftCluster) error {
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LenshoodRaftClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&LenshoodRaftCluster{}).
		Complete(r)
}
