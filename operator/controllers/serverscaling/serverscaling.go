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

package serverscaling

import (
	"context"
	"fmt"

	mlopsv1alpha1 "github.com/seldonio/seldon-core/operator/v2/apis/mlops/v1alpha1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ServerScaleEventMsg struct {
	Namespace  string
	ServerName string
	Replicas   int32
}

type coordinator struct {
	events chan ServerScaleEventMsg
	client client.Client
}

func NewCoordinator(events chan ServerScaleEventMsg, client client.Client) coordinator {
	return coordinator{
		events: events,
		client: client,
	}
}

func (c *coordinator) StartWatchServerScaleEvents() {
	go func() {
		for e := range c.events {
			c.handleServerScaleEvent(e)
		}
	}()
}

func (c *coordinator) handleServerScaleEvent(event ServerScaleEventMsg) {
	logger := log.FromContext(context.TODO()).WithName("handleServerScaleEvent")
	var err error
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		server := &mlopsv1alpha1.Server{}
		if err = c.client.Get(context.TODO(), client.ObjectKey{Name: event.ServerName, Namespace: event.Namespace}, server); err != nil {
			logger.Error(err, "unable to get server", "name", event.ServerName, "namespace", event.Namespace)
			return err
		}
		if server.Spec.Replicas != nil && *server.Spec.Replicas == event.Replicas {
			return nil
		}
		replicas := event.Replicas
		server.Spec.ScalingSpec.Replicas = &replicas

		err = c.client.Update(context.TODO(), server)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("scaled replicas to %d for server %s", event.Replicas, server.Name))
		return nil
	})
	if retryErr != nil {
		logger.Error(err, fmt.Sprintf("failed to scale replicas to %d for server %s", event.Replicas, event.ServerName))
	}
}
