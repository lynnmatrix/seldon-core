/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package filters

import (
	"testing"

	. "github.com/onsi/gomega"

	pb "github.com/seldonio/seldon-core/apis/go/v2/mlops/scheduler"

	"github.com/seldonio/seldon-core/scheduler/v2/pkg/store"
)

func getTestModelWithMemory(requiredmemory *uint64, serverName string, replicaId int) *store.ModelVersion {

	replicas := map[int]store.ReplicaStatus{}
	if replicaId >= 0 {
		replicas[replicaId] = store.ReplicaStatus{State: store.Loading}
	}
	return store.NewModelVersion(
		&pb.Model{ModelSpec: &pb.ModelSpec{MemoryBytes: requiredmemory}, DeploymentSpec: &pb.DeploymentSpec{Replicas: 1}},
		1,
		serverName,
		replicas,
		false,
		store.ModelProgressing)
}

func getTestServerReplicaWithMemory(availableMemory, reservedMemory uint64, serverName string, replicaId int) *store.ServerReplica {
	return store.NewServerReplica("svc", 8080, 5001, replicaId, store.NewServer(serverName, true), []string{}, availableMemory, availableMemory, reservedMemory, nil, 100)
}

func TestReplicaMemoryFilter(t *testing.T) {
	g := NewGomegaWithT(t)

	type test struct {
		name     string
		affinity bool
		model    *store.ModelVersion
		server   *store.ServerReplica
		expected bool
	}

	memory := uint64(100)
	tests := []test{
		{name: "EnoughMemory", affinity: true, model: getTestModelWithMemory(&memory, "", -1), server: getTestServerReplicaWithMemory(100, 0, "server1", 0), expected: true},
		{name: "NoMemorySpecified", affinity: true, model: getTestModelWithMemory(nil, "", -1), server: getTestServerReplicaWithMemory(200, 0, "server1", 0), expected: true},
		{name: "NotEnoughMemory", affinity: true, model: getTestModelWithMemory(&memory, "", -1), server: getTestServerReplicaWithMemory(50, 0, "server1", 0), expected: false},
		{name: "NotEnoughMemoryWithReserved", affinity: true, model: getTestModelWithMemory(&memory, "", -1), server: getTestServerReplicaWithMemory(200, 150, "server1", 0), expected: false},
		{name: "NotEnoughMemoryWithReservedOverflow", affinity: true, model: getTestModelWithMemory(&memory, "", -1), server: getTestServerReplicaWithMemory(200, 250, "server1", 0), expected: false},
		{name: "ModelAlreadyLoaded", affinity: true, model: getTestModelWithMemory(&memory, "server1", 0), server: getTestServerReplicaWithMemory(0, 0, "server1", 0), expected: true}, // note not enough memory on server replica
		{name: "AntiAffiniy", affinity: false, model: getTestModelWithMemory(&memory, "server1", 0), server: getTestServerReplicaWithMemory(200, 0, "server1", 0), expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter := AvailableMemoryReplicaFilter{Affinity: test.affinity}
			ok := filter.Filter(test.model, test.server)
			g.Expect(ok).To(Equal(test.expected))
		})
	}
}
