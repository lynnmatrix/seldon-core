/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package scheduler

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/seldonio/seldon-core/scheduler/v2/pkg/scheduler/filters"
	"github.com/seldonio/seldon-core/scheduler/v2/pkg/scheduler/sorters"
	store "github.com/seldonio/seldon-core/scheduler/v2/pkg/store"
)

type SimpleScheduler struct {
	muSortAndUpdate sync.Mutex
	store           store.ModelStore
	logger          log.FieldLogger
	SchedulerConfig
}

type SchedulerConfig struct {
	serverFilters        []filters.ServerFilter
	serverSorts          []sorters.ServerSorter
	replicaFilters       []filters.ReplicaFilter
	replicaSorts         []sorters.ReplicaSorter
	serverScalableFilter []filters.ReplicaFilter
}

func DefaultSchedulerConfig(store store.ModelStore) SchedulerConfig {
	return SchedulerConfig{
		serverFilters:  []filters.ServerFilter{filters.ServerReplicaFilter{}, filters.SharingServerFilter{}, filters.DeletedServerFilter{}, filters.ServerRequirementFilter{}},
		replicaFilters: []filters.ReplicaFilter{filters.AvailableMemoryReplicaFilter{}, filters.ExplainerFilter{}, filters.ReplicaDrainingFilter{}},
		serverSorts:    []sorters.ServerSorter{},
		replicaSorts:   []sorters.ReplicaSorter{sorters.ReplicaIndexSorter{}, sorters.AvailableMemorySorter{}, sorters.ModelAlreadyLoadedSorter{}},
		serverScalableFilter: []filters.ReplicaFilter{filters.ExplainerFilter{}},
	}
}

func NewSimpleScheduler(logger log.FieldLogger,
	store store.ModelStore,
	schedulerConfig SchedulerConfig) *SimpleScheduler {
	s := &SimpleScheduler{
		store:           store,
		logger:          logger.WithField("Name", "SimpleScheduler"),
		SchedulerConfig: schedulerConfig,
	}
	return s
}

func (s *SimpleScheduler) Schedule(modelKey string) error {
	return s.scheduleToServer(modelKey)
}

func (s *SimpleScheduler) ScheduleFailedModels() ([]string, error) {
	failedModels, err := s.getFailedModels()
	if err != nil {
		return nil, err
	}
	var updatedModels []string
	for _, modelName := range failedModels {
		err := s.scheduleToServer(modelName)
		if err != nil {
			s.logger.Debugf("Failed to schedule failed model %s", modelName)
		} else {
			updatedModels = append(updatedModels, modelName)
		}
	}
	return updatedModels, nil
}

func (s *SimpleScheduler) getFailedModels() ([]string, error) {
	models, err := s.store.GetModels()
	if err != nil {
		return nil, err
	}
	var failedModels []string
	for _, model := range models {
		version := model.GetLatest()
		if version != nil {
			versionState := version.ModelState()
			if versionState.State == store.ModelFailed || versionState.State == store.ScheduleFailed {
				failedModels = append(failedModels, model.Name)
			}
		}
	}
	return failedModels, nil
}

// TODO - clarify non shared models should not be scheduled
func (s *SimpleScheduler) scheduleToServer(modelName string) error {
	logger := s.logger.WithField("func", "scheduleToServer").WithField("model", modelName)
	logger.Debug("Schedule model")

	s.store.LockModel(modelName)
	defer s.store.UnlockModel(modelName)

	// Get Model
	model, err := s.store.GetModel(modelName)
	if err != nil {
		return err
	}
	if model == nil {
		return errors.New("Unable to find model")
	}

	latestModel := model.GetLatest()
	if latestModel == nil {
		return errors.New("Unable to find latest version for model")
	}

	if model.Deleted {
		// we need to LoadedModels anyway:
		// - in case where we are deleting a model that doesnt have a server (FailedSchedule), server is ""
		// - otherwise proceed a normal
		server := ""
		if latestModel.HasServer() {
			server = latestModel.Server()
		}

		logger.Debug("Ensuring deleted model is removed")
		err = s.store.UpdateLoadedModels(modelName, latestModel.GetVersion(), server, []*store.ServerReplica{})
		if err != nil {
			logger.WithError(err).WithField("server", server).Warn("Failed to unschedule model replicas from server")
		}
		return nil
	} else {
		// Model needs to be (re)scheduled
		var filteredServers []*store.ServerSnapshot
		
		// Get all servers
		servers, err := s.store.GetServers(false, true)
		if err != nil {
			return err
		}
		// Filter and sort servers
		filteredServers = s.filterServers(latestModel, servers)
		if len(filteredServers) == 0 {
			msg := "Failed to schedule model as no matching servers are available"
			logger.Debug(msg)
			s.store.FailedScheduling(latestModel, msg, !latestModel.HasLiveReplicas())
			return errors.New(msg)
		}

		s.sortServers(latestModel, filteredServers)
		logger.
			WithField("candidate_servers", filteredServers).
			WithField("desired_replicas", latestModel.DesiredReplicas()).
			Debug("Identified candidate servers for model")

		// For each server filter and sort replicas and attempt schedule if enough replicas
		ok := false
		for _, candidateServer := range filteredServers {
			logger.WithField("server", candidateServer.Name).Debug("Checking compatibility with candidate server")
			var candidateReplicas *sorters.CandidateServer

			// we need a lock here, we could have many goroutines at sorting
			// without the store being reflected and hence sorting on stale values
			s.muSortAndUpdate.Lock()
			candidateReplicas = s.filterReplicas(latestModel, candidateServer)
			if len(candidateReplicas.ChosenReplicas) < latestModel.DesiredReplicas() {
				
				s.muSortAndUpdate.Unlock()

				scaleToReplicas := calScaleToReplicas(candidateServer, latestModel.DesiredReplicas(), len(candidateReplicas.ChosenReplicas))
				if s.serverScalable(candidateServer, latestModel, scaleToReplicas) {
					logger.Debugf("scale up server %s to %d replicas for model %s", candidateServer.Name, scaleToReplicas, latestModel.GetMeta().Name)
					s.store.UpdateServerScaleToReplicas(candidateServer.Name, int32(scaleToReplicas))
					break
				}
				logger.Debugf("cann't scale up server %s to %d replicas for model %s", candidateServer.Name, scaleToReplicas, latestModel.GetMeta().Name)
				
				continue
			}
			s.sortReplicas(candidateReplicas)
			err = s.store.UpdateLoadedModels(
				modelName, latestModel.GetVersion(), 
				candidateServer.Name, 
				candidateReplicas.ChosenReplicas[0:latestModel.DesiredReplicas()],
			)
			s.muSortAndUpdate.Unlock()

			if err != nil {
				logger.WithField("server", candidateServer.Name).Warn("Failed to update model replicas")
			} else {
				logger.WithField("server", candidateServer.Name).Debug("Scheduled model onto server")
				ok = true
				break
			}
		}
	
		if !ok {
			msg := "Failed to schedule model as no matching server had enough suitable replicas"
			logger.Debug(msg)
			// we do not want to reset the server if it has live replicas
			s.store.FailedScheduling(latestModel, msg, !latestModel.HasLiveReplicas())
			return errors.New(msg)
		}
	}

	//TODO Cleanup previous version if needed?

	return s.scaleDownServerIfNeed()
}

// FIXME 暂时采用最简单的策略，当最大的模型副本数小于 server 副本数时，触发scale down。
// 这个策略没有考虑到 available memory, 因此会出现scale down 过于激进的问题，需要尽快完善
func (s *SimpleScheduler) scaleDownServerIfNeed() error {
	serverMaxModelReplicas := map[string]int{}
	models, err := s.store.GetModels()
	if err != nil {
		return err
	}
	for _, model := range models {
		modelVersion := model.GetLatest()
		if !modelVersion.HasServer() {
			continue
		}
		serverKey := modelVersion.Server()
		maxModelReplicas, ok := serverMaxModelReplicas[modelVersion.Server()]
		if !ok {
			maxModelReplicas = 0
		}

		if modelVersion.DesiredReplicas() > maxModelReplicas {
			serverMaxModelReplicas[serverKey] = modelVersion.DesiredReplicas()
		}
	}

	for serverKey, maxModelReplicas := range serverMaxModelReplicas {
		server, err := s.store.GetServer(serverKey, true, false)
		if err != nil {
			return err
		}
		if server.ExpectedReplicas > maxModelReplicas {
			scaleToReplicas := calScaleToReplicas(server, maxModelReplicas, server.ExpectedReplicas)
			if s.serverScalable(server, nil, scaleToReplicas) {
				s.logger.Debugf("scale down server %s to %d replicas", server.Name, scaleToReplicas)
				s.store.UpdateServerScaleToReplicas(server.Name, int32(scaleToReplicas))
			}
		}
	}
	return nil
}

func calScaleToReplicas(server *store.ServerSnapshot, desiredReplicas int, availableReplicas int) int {
	scaleToReplicas := server.ExpectedReplicas + (desiredReplicas - availableReplicas)
	if server.MaxReplicas > 0 && scaleToReplicas > server.MaxReplicas {
		scaleToReplicas = server.MaxReplicas
	}
	if server.MinReplicas > 0 && scaleToReplicas < server.MinReplicas {
		scaleToReplicas = server.MinReplicas
	}
	return scaleToReplicas
}

func (s *SimpleScheduler) serverScalable(server *store.ServerSnapshot, model *store.ModelVersion, scaleToReplicas int) bool {

	logger := s.logger.WithField("func", "serverScalable")

	autoscalingEnabled := server.MinReplicas > 0 || server.MaxReplicas > 0
	if !autoscalingEnabled {
		return false
	}

	// check whether the server is during scaling
	if len(server.Replicas) != server.ExpectedReplicas {
		logger.Debugf("server %s is not scalable as its state is not stable", server.Name)
		return false
	}

	// check the capbilities of server
	if model != nil {
		capable := false
		for _, replica := range server.Replicas {
			ok := true
			for _, replicaFilter := range s.serverScalableFilter {
				if !replicaFilter.Filter(model, replica) {
					ok = false
					break
				}
			}
			if ok {
				capable = true
				break
			}
		}

		if !capable {
			logger.Debugf("server %s is not scalable for model %s because of capability mismatch", server.Name, model.GetMeta().Name)
			return false
		}
	}

	return checkDesiredNumReplicas(server, scaleToReplicas)
}

func checkDesiredNumReplicas(server *store.ServerSnapshot, replicas int) bool {
	if replicas == server.ExpectedReplicas {
		return false
	} else if replicas > server.ExpectedReplicas {
		return replicas <= server.MaxReplicas || server.MaxReplicas <= 0
	} else {
		return replicas >= server.MinReplicas && replicas > 0
	}
}

func showServerSlice(servers []*store.ServerSnapshot) string {
	var sb strings.Builder
	for idx, server := range servers {
		if idx > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(server.Name)
	}
	return sb.String()
}

func (s *SimpleScheduler) sortServers(model *store.ModelVersion, server []*store.ServerSnapshot) {
	logger := s.logger.WithField("func", "sortServers")
	for _, sorter := range s.serverSorts {
		logger.Debugf("About to sort servers for %s:%d with %s: %s", model.Key(), model.GetVersion(), sorter.Name(), showServerSlice(server))
		sort.SliceStable(server, func(i, j int) bool {
			return sorter.IsLess(&sorters.CandidateServer{Model: model, Server: server[i]}, &sorters.CandidateServer{Model: model, Server: server[j]})
		})
		logger.Debugf("Sorted servers for %s:%d with %s: %s", model.Key(), model.GetVersion(), sorter.Name(), showServerSlice(server))
	}
}

func showReplicaSlice(candidateServer *sorters.CandidateServer) string {
	var sb strings.Builder
	for idx, replica := range candidateServer.ChosenReplicas {
		if idx > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.Itoa(replica.GetReplicaIdx()))
		sb.WriteString(":")
		sb.WriteString(replica.GetInferenceSvc())
	}
	return sb.String()
}

func (s *SimpleScheduler) sortReplicas(candidateServer *sorters.CandidateServer) {
	logger := s.logger.WithField("func", "sortReplicas")
	for _, sorter := range s.replicaSorts {
		logger.Debugf("About to sort replicas for %s:%d with %s: %s", candidateServer.Model.Key(), candidateServer.Model.GetVersion(), sorter.Name(), showReplicaSlice(candidateServer))
		sort.SliceStable(candidateServer.ChosenReplicas, func(i, j int) bool {
			return sorter.IsLess(&sorters.CandidateReplica{Model: candidateServer.Model, Server: candidateServer.Server, Replica: candidateServer.ChosenReplicas[i]},
				&sorters.CandidateReplica{Model: candidateServer.Model, Server: candidateServer.Server, Replica: candidateServer.ChosenReplicas[j]})
		})
		logger.Debugf("Sorted replicas for %s:%d with %s: %s", candidateServer.Model.Key(), candidateServer.Model.GetVersion(), sorter.Name(), showReplicaSlice(candidateServer))
	}
}

// Filter servers for this model
func (s *SimpleScheduler) filterServers(model *store.ModelVersion, servers []*store.ServerSnapshot) []*store.ServerSnapshot {
	logger := s.logger.WithField("func", "filterServer").WithField("model", model.GetMeta().GetName())
	logger.WithField("num_servers", len(servers)).Debug("Filtering servers for model")

	var filteredServers []*store.ServerSnapshot
	for _, server := range servers {
		ok := true
		for _, serverFilter := range s.serverFilters {
			if !serverFilter.Filter(model, server) {
				logger.
					WithField("filter", serverFilter.Name()).
					WithField("server", server.Name).
					WithField("reason", serverFilter.Description(model, server)).
					Debug("Rejecting server for model")

				ok = false
				break
			}
		}

		if ok {
			logger.WithField("server", server.Name).Debug("Accepting server for model")
			filteredServers = append(filteredServers, server)
		}
	}

	return filteredServers
}

func (s *SimpleScheduler) filterReplicas(model *store.ModelVersion, server *store.ServerSnapshot) *sorters.CandidateServer {
	logger := s.logger.
		WithField("func", "filterReplicas").
		WithField("model", model.GetMeta().GetName()).
		WithField("server", server.Name)
	logger.Debug("Filtering server replicas for model")

	candidateServer := sorters.CandidateServer{Model: model, Server: server}
	for _, replica := range server.Replicas {
		ok := true
		for _, replicaFilter := range s.replicaFilters {
			if !replicaFilter.Filter(model, replica) {
				logger.
					WithField("filter", replicaFilter.Name()).
					WithField("replica", replica.GetReplicaIdx()).
					WithField("reason", replicaFilter.Description(model, replica)).
					Debug("Rejecting server replica for model")

				ok = false
				break
			}
		}

		if ok {
			logger.WithField("replica", replica.GetReplicaIdx()).Debug("Accepting server replica for model")
			candidateServer.ChosenReplicas = append(candidateServer.ChosenReplicas, replica)
		}
	}

	return &candidateServer
}
