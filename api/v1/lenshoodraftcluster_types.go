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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LenshoodRaftClusterSpec defines the desired state of LenshoodRaftCluster
type LenshoodRaftClusterSpec struct {
	// Replica defines how many raft instance exists in a single cluster
	Replica int `json:"replica"`

	// Image defines lenshood raft image
	Image string `json:"image"`
}

type ClusterState uint8

const (
	Init ClusterState = iota
	OK
	BUILDING
	ERROR
	UNKNOWN
)

// LenshoodRaftClusterStatus defines the observed state of LenshoodRaftCluster
type LenshoodRaftClusterStatus struct {
	State ClusterState `json:"state"`
	Ids   []string     `json:"ids"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LenshoodRaftCluster is the Schema for the lenshoodraftclusters API
type LenshoodRaftCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LenshoodRaftClusterSpec   `json:"spec,omitempty"`
	Status LenshoodRaftClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LenshoodRaftClusterList contains a list of LenshoodRaftCluster
type LenshoodRaftClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LenshoodRaftCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LenshoodRaftCluster{}, &LenshoodRaftClusterList{})
}
