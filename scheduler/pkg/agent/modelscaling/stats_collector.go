/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package modelscaling

import (
	"github.com/seldonio/seldon-core/scheduler/v2/pkg/agent/interfaces"
)

type DataPlaneStatsCollector struct {
	StatKeepers []interfaces.ModelStatsKeeper
}

func NewDataPlaneStatsCollector(statKeepers []interfaces.ModelStatsKeeper) *DataPlaneStatsCollector {
	return &DataPlaneStatsCollector{
		StatKeepers: statKeepers,
	}
}

func (c *DataPlaneStatsCollector) ModelInferEnter(internalModelName, requestId string) error {
	var err error
	for _, stat := range c.StatKeepers {
		err = stat.ModelInferEnter(internalModelName, requestId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *DataPlaneStatsCollector) ModelInferExit(internalModelName, requestId string) error {
	var err error
	for _, stat := range c.StatKeepers {
		err = stat.ModelInferExit(internalModelName, requestId)
		if err != nil {
			return err
		}
	}
	return nil
}
