/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package scheduler

import "github.com/seldonio/seldon-core/scheduler/v2/pkg/store"

type Scheduler interface {
	Schedule(modelKey string) error
	ScheduleFailedModels() ([]string, error)
}

type ServerScaler interface {
	Scalable(serverKey string, replicas int, model *store.ModelVersion) bool
}
