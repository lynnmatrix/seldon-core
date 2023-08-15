/*
Copyright 2023 Seldon Technologies Ltd.

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

package filters

import (
	"fmt"

	"github.com/seldonio/seldon-core/scheduler/v2/pkg/store"
)

type MaxIndexReplicaFilter struct {
	MaxIndex uint
}

func (r MaxIndexReplicaFilter) Name() string {
	return "MaxIndexReplicaFilter"
}

func (r MaxIndexReplicaFilter) Filter(model *store.ModelVersion, replica *store.ServerReplica) bool {
	return replica.GetReplicaIdx() <= int(r.MaxIndex)
}

func (r MaxIndexReplicaFilter) Description(model *store.ModelVersion, replica *store.ServerReplica) string {
	return fmt.Sprintf("max idx %d replica idx %d", r.MaxIndex, replica.GetReplicaIdx())
}
