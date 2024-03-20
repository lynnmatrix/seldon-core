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
	scaler          ServerScaler
	SchedulerConfig
}

type SchedulerConfig struct {
	serverFilters  []filters.ServerFilter
	serverSorts    []sorters.ServerSorter
	replicaFilters []filters.ReplicaFilter
	replicaSorts   []sorters.ReplicaSorter
}

func DefaultSchedulerConfig(store store.ModelStore) SchedulerConfig {
	return SchedulerConfig{
		serverFilters:  []filters.ServerFilter{filters.ServerReplicaFilter{}, filters.SharingServerFilter{}, filters.DeletedServerFilter{}, filters.ServerRequirementFilter{}},
		replicaFilters: []filters.ReplicaFilter{filters.AvailableMemoryReplicaFilter{}, filters.ExplainerFilter{}, filters.ReplicaDrainingFilter{}},
		serverSorts:    []sorters.ServerSorter{},
		replicaSorts:   []sorters.ReplicaSorter{sorters.ReplicaIndexSorter{}, sorters.AvailableMemorySorter{}, sorters.ModelAlreadyLoadedSorter{}},
	}
}

func NewSimpleScheduler(logger log.FieldLogger,
	store store.ModelStore,
	schedulerConfig SchedulerConfig,
	scaler ServerScaler) *SimpleScheduler {
	s := &SimpleScheduler{
		store:           store,
		logger:          logger.WithField("Name", "SimpleScheduler"),
		scaler:          scaler,
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
			candidateReplicas = filterReplicas(s.replicaFilters, latestModel, candidateServer, s.logger)

			if len(candidateReplicas.ChosenReplicas) < latestModel.DesiredReplicas() {
				scaleToReplicas := calScaleToReplicas(candidateServer, latestModel.DesiredReplicas(), len(candidateReplicas.ChosenReplicas))

				if s.scaler.Scalable(candidateServer.Name, scaleToReplicas, latestModel) {
					logger.Debugf("scale up server %s to %d replicas for model %s", candidateServer.Name, scaleToReplicas, modelName)
					s.store.UpdateServerScaleToReplicas(candidateServer.Name, int32(scaleToReplicas))
					s.muSortAndUpdate.Unlock()
					break
				}

				s.muSortAndUpdate.Unlock()
				logger.Debugf("cann't scale up server %s to %d replicas for model %s", candidateServer.Name, scaleToReplicas, latestModel.GetMeta().Name)
				
				continue
			}
			candidateReplicas.SortReplicas(s.replicaSorts)
			err = s.store.UpdateLoadedModels(
				modelName, latestModel.GetVersion(), 
				candidateServer.Name, 
				candidateReplicas.ChosenReplicas[0:latestModel.DesiredReplicas()])
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

func (s *SimpleScheduler) scaleDownServerIfNeed() error {
	servers, err := s.store.GetServers(false, true)
	if err != nil {
		return err
	}

	for _, server := range servers {
		scaleToReplicas := server.ExpectedReplicas - 1

		if s.scaler.Scalable(server.Name, scaleToReplicas, nil) {
			s.logger.Debugf("scale down server %s to %d replicas", server.Name, scaleToReplicas)
			s.store.UpdateServerScaleToReplicas(server.Name, int32(scaleToReplicas))
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
