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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
			return res, r.deleteCluster(ctx, cluster.DeepCopy())
		}
		return
	}

	if r.reachedDesiredState(ctx, cluster) {
		return
	}

	if len(cluster.Status.Ids) == 0 {
		return res, r.createSeedToCluster(ctx, cluster.DeepCopy())
	}

	return ctrl.Result{}, r.joinToCluster(ctx, cluster.DeepCopy())
}

func (r *LenshoodRaftClusterReconciler) reachedDesiredState(ctx context.Context, cluster *LenshoodRaftCluster) bool {
	if cluster.Status.State != OK {
		return false
	}

	if len(cluster.Status.Ids) != cluster.Spec.Replica {
		return false
	}

	for _, id := range cluster.Status.Ids {
		if err := r.Client.Get(ctx, types.NamespacedName{Namespace: id, Name: cluster.Namespace}, &v1.Pod{}); err != nil {
			return false
		}
	}

	return true
}

func (r *LenshoodRaftClusterReconciler) createSeedToCluster(ctx context.Context, cluster *LenshoodRaftCluster) error {
	return nil
}

func (r *LenshoodRaftClusterReconciler) joinToCluster(ctx context.Context, cluster *LenshoodRaftCluster) error {
	return nil
}

func (r *LenshoodRaftClusterReconciler) deleteCluster(ctx context.Context, cluster *LenshoodRaftCluster) error {
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LenshoodRaftClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&LenshoodRaftCluster{}).
		Complete(r)
}
