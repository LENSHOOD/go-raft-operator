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
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog/v2"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"

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
		return res, r.createCluster(ctx, cluster.DeepCopy())
	}

	return ctrl.Result{}, r.updateCluster(ctx, cluster.DeepCopy())
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

func getPodName(ordinal int) string {
	return fmt.Sprintf("lenshood-raft-%d", ordinal)
}

func buildRaftHeadlessSvc(cluster *LenshoodRaftCluster) *v1.Service {
	svcYaml, err := os.ReadFile("config/svc/raft-svc.yaml")
	if err != nil {
		klog.Fatalf("read pod yaml file failed: %v", err)
	}

	svc := &v1.Service{}
	if err := yaml.Unmarshal(svcYaml, svc); err != nil {
		klog.Fatalf("unmarshal pod yaml file failed: %v", err)
	}

	svc.Namespace = cluster.Namespace
	svc.Name = cluster.Name
	svc.Spec.Selector = map[string]string{"managed-by": cluster.Name}
	return svc
}

func buildRaftPod(cluster *LenshoodRaftCluster, members []string, ordinal int) *v1.Pod {
	podYaml, err := os.ReadFile("config/pod/raft-pod.yaml")
	if err != nil {
		klog.Fatalf("read pod yaml file failed: %v", err)
	}

	pod := &v1.Pod{}
	if err := yaml.Unmarshal(podYaml, pod); err != nil {
		klog.Fatalf("unmarshal pod yaml file failed: %v", err)
	}

	pod.ObjectMeta.Labels = map[string]string{"managed-by": cluster.Name}
	pod.Namespace = cluster.Namespace
	pod.Name = getPodName(ordinal)
	pod.Spec.Hostname = getPodName(ordinal)
	pod.Spec.Subdomain = cluster.Name
	pod.Spec.Containers[0].Image = cluster.Spec.Image
	pod.Spec.Containers[0].Env = []v1.EnvVar{{Name: "ME_ADDR", Value: members[ordinal]}, {Name: "MEMBERS_ADDR", Value: strings.Join(members, ",")}}

	return pod
}

func (r *LenshoodRaftClusterReconciler) createCluster(ctx context.Context, cluster *LenshoodRaftCluster) error {
	svc := buildRaftHeadlessSvc(cluster)
	if err := r.Create(ctx, svc); err != nil {
		return err
	}

	podNums := cluster.Spec.Replica
	members := make([]string, 0, podNums)
	for i := 0; i < podNums; i++ {
		members = append(members, getPodName(i)+"."+svc.Name+":34220")
	}

	for i := 0; i < cluster.Spec.Replica; i++ {
		pod := buildRaftPod(cluster, members, i)
		if err := r.Create(ctx, pod); err != nil {
			return err
		}
		cluster.Status.Ids = append(cluster.Status.Ids, pod.Name)
	}

	return nil
}

func (r *LenshoodRaftClusterReconciler) updateCluster(ctx context.Context, cluster *LenshoodRaftCluster) error {
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
