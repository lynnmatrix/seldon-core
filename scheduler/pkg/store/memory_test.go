package store

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	pb "github.com/seldonio/seldon-core/scheduler/apis/mlops/scheduler"
	log "github.com/sirupsen/logrus"
)

func TestUpdateModel(t *testing.T) {
	logger := log.New()
	g := NewGomegaWithT(t)

	type test struct {
		name   string
		store  *LocalSchedulerStore
		config *pb.ModelDetails
		err    error
	}

	tests := []test{
		{
			name:   "simple",
			store:  NewLocalSchedulerStore(),
			config: &pb.ModelDetails{Name: "model", Version: "1"},
			err:    nil,
		},
		{
			name: "VersionAlreadyExists",
			store: &LocalSchedulerStore{models: map[string]*Model{"model": &Model{
				versionMap: map[string]*ModelVersion{"1": &ModelVersion{
					config: &pb.ModelDetails{Name: "model", Version: "1"},
				}},
				versions: []*ModelVersion{&ModelVersion{config: &pb.ModelDetails{Name: "model", Version: "1"}}},
			}}},
			config: &pb.ModelDetails{Name: "model", Version: "1"},
			err:    ModelVersionExistsErr,
		},
		{
			name: "VersionAdded",
			store: &LocalSchedulerStore{models: map[string]*Model{"model": &Model{
				versionMap: map[string]*ModelVersion{"1": &ModelVersion{
					config: &pb.ModelDetails{Name: "model", Version: "1"},
				}},
				versions: []*ModelVersion{&ModelVersion{config: &pb.ModelDetails{Name: "model", Version: "1"}}},
			}}},
			config: &pb.ModelDetails{Name: "model", Version: "2"},
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := NewMemoryStore(logger, test.store)
			err := ms.UpdateModel(test.config)
			if test.err == nil {
				g.Expect(err).To(BeNil())
				m := test.store.models[test.config.Name]
				g.Expect(m.versions[len(m.versions)-1].config).To(Equal(test.config))
			} else {
				g.Expect(err).ToNot(BeNil())
			}
		})
	}
}

func TestGetModel(t *testing.T) {
	logger := log.New()
	g := NewGomegaWithT(t)

	type test struct {
		name     string
		store    *LocalSchedulerStore
		key      string
		versions int
		err      error
	}

	tests := []test{
		{
			name:     "NoModel",
			store:    NewLocalSchedulerStore(),
			key:      "model",
			versions: 0,
			err:      nil,
		},
		{
			name: "VersionAlreadyExists",
			store: &LocalSchedulerStore{models: map[string]*Model{"model": &Model{
				versionMap: map[string]*ModelVersion{"1": &ModelVersion{
					config: &pb.ModelDetails{Name: "model", Version: "1"},
				}},
				versions: []*ModelVersion{&ModelVersion{config: &pb.ModelDetails{Name: "model", Version: "1"}}},
			}}},
			key:      "model",
			versions: 1,
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := NewMemoryStore(logger, test.store)
			model, err := ms.GetModel(test.key)
			if test.err == nil {
				g.Expect(err).To(BeNil())
				g.Expect(model.Name).To(Equal(test.key))
				g.Expect(len(model.Versions)).To(Equal(test.versions))
			} else {
				g.Expect(err).ToNot(BeNil())
			}
		})
	}
}

func TestRemoveModel(t *testing.T) {
	logger := log.New()
	g := NewGomegaWithT(t)

	type test struct {
		name  string
		store *LocalSchedulerStore
		key   string
		err   error
	}

	tests := []test{
		{
			name:  "NoModel",
			store: NewLocalSchedulerStore(),
			key:   "model",
			err:   nil,
		},
		{
			name: "VersionAlreadyExists",
			store: &LocalSchedulerStore{models: map[string]*Model{"model": &Model{
				versionMap: map[string]*ModelVersion{"1": &ModelVersion{
					config: &pb.ModelDetails{Name: "model", Version: "1"},
				}},
				versions: []*ModelVersion{&ModelVersion{config: &pb.ModelDetails{Name: "model", Version: "1"}}},
			}}},
			key: "model",
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := NewMemoryStore(logger, test.store)
			err := ms.RemoveModel(test.key)
			if test.err == nil {
				g.Expect(err).To(BeNil())
			} else {
				g.Expect(err).ToNot(BeNil())
			}
		})
	}
}

func TestUpdateLoadedModels(t *testing.T) {
	logger := log.New()
	g := NewGomegaWithT(t)

	type test struct {
		name           string
		store          *LocalSchedulerStore
		modelKey       string
		version        string
		serverKey      string
		replicas       []*ServerReplica
		expectedStates map[int]ModelReplicaState
		err            error
	}

	tests := []test{
		{
			name: "ModelVersionNotLatest",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config:   &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{},
						},
						{
							config:   &pb.ModelDetails{Name: "model", Version: "2"},
							replicas: map[int]ModelReplicaState{},
						},
					},
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
					},
				},
			},
			modelKey:  "model",
			version:   "1",
			serverKey: "server",
			replicas:  nil,
			err:       ModelNotLatestVersionRejectErr,
		},
		{
			name: "UpdatedVersionsOK",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config:   &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{},
						},
					},
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
						replicas: map[int]*ServerReplica{
							0: {},
							1: {},
						},
					},
				},
			},
			modelKey:  "model",
			version:   "1",
			serverKey: "server",
			replicas: []*ServerReplica{
				{replicaIdx: 0}, {replicaIdx: 1},
			},
			expectedStates: map[int]ModelReplicaState{0: LoadRequested, 1: LoadRequested},
			err:            nil,
		},
		{
			name: "WithAlreadyLoadedModels",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config: &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{
								0: Loaded,
							},
						},
					},
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
						replicas: map[int]*ServerReplica{
							0: {},
							1: {},
						},
					},
				},
			},
			modelKey:  "model",
			version:   "1",
			serverKey: "server",
			replicas: []*ServerReplica{
				{replicaIdx: 0}, {replicaIdx: 1},
			},
			expectedStates: map[int]ModelReplicaState{0: Loaded, 1: LoadRequested},
			err:            nil,
		},
		{
			name: "UnloadModelsNotSelected",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config: &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{
								0: Loaded,
							},
						},
					},
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
						replicas: map[int]*ServerReplica{
							0: {},
							1: {},
						},
					},
				},
			},
			modelKey:  "model",
			version:   "1",
			serverKey: "server",
			replicas: []*ServerReplica{
				{replicaIdx: 1},
			},
			expectedStates: map[int]ModelReplicaState{0: UnloadRequested, 1: LoadRequested},
			err:            nil,
		},
		{
			name: "DeletedModel",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config:   &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{0: Loaded},
						},
					},
					deleted: true,
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
						replicas: map[int]*ServerReplica{
							0: {},
							1: {},
						},
					},
				},
			},
			modelKey:       "model",
			version:        "1",
			serverKey:      "server",
			replicas:       []*ServerReplica{},
			expectedStates: map[int]ModelReplicaState{0: UnloadRequested},
			err:            nil,
		},
		{
			name: "DeletedModelNoReplicas",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config:   &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{0: Unloaded},
						},
					},
					deleted: true,
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
						replicas: map[int]*ServerReplica{
							0: {},
							1: {},
						},
					},
				},
			},
			modelKey:       "model",
			version:        "1",
			serverKey:      "server",
			replicas:       []*ServerReplica{},
			expectedStates: map[int]ModelReplicaState{0: Unloaded},
			err:            nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := NewMemoryStore(logger, test.store)
			err := ms.UpdateLoadedModels(test.modelKey, test.version, test.serverKey, test.replicas)
			if test.err == nil {
				g.Expect(err).To(BeNil())
				for replicaIdx, state := range test.expectedStates {
					g.Expect(test.store.models[test.modelKey].Latest().GetModelReplicaState(replicaIdx)).To(Equal(state))
				}
			} else {
				g.Expect(err).ToNot(BeNil())
				g.Expect(errors.Is(err, test.err)).To(BeTrue())
			}
		})
	}
}

func TestUpdateModelState(t *testing.T) {
	logger := log.New()
	g := NewGomegaWithT(t)

	type test struct {
		name            string
		store           *LocalSchedulerStore
		modelKey        string
		version         string
		serverKey       string
		replicaIdx      int
		state           ModelReplicaState
		availableMemory uint64
		loaded          bool
		deleted         bool
		err             error
	}

	tests := []test{
		{
			name: "LoadedModel",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config:   &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{},
						},
					},
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
						replicas: map[int]*ServerReplica{
							0: {loadedModels: map[string]bool{}},
							1: {loadedModels: map[string]bool{}},
						},
					},
				},
			},
			modelKey:        "model",
			version:         "1",
			serverKey:       "server",
			replicaIdx:      0,
			state:           Loaded,
			loaded:          true,
			availableMemory: 20,
			err:             nil,
		},
		{
			name: "UnloadedModel",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config:   &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{},
						},
					},
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
						replicas: map[int]*ServerReplica{
							0: {loadedModels: map[string]bool{}},
							1: {loadedModels: map[string]bool{}},
						},
					},
				},
			},
			modelKey:        "model",
			version:         "1",
			serverKey:       "server",
			replicaIdx:      0,
			state:           Unloaded,
			loaded:          false,
			availableMemory: 20,
			err:             nil,
		},
		{
			name: "DeletedModel",
			store: &LocalSchedulerStore{
				models: map[string]*Model{"model": &Model{
					versionMap: map[string]*ModelVersion{"1": {
						config:   &pb.ModelDetails{Name: "model", Version: "1"},
						replicas: map[int]ModelReplicaState{},
					}},
					versions: []*ModelVersion{
						{
							config:   &pb.ModelDetails{Name: "model", Version: "1"},
							replicas: map[int]ModelReplicaState{},
						},
					},
					deleted: true,
				}},
				servers: map[string]*Server{
					"server": {
						name: "server",
						replicas: map[int]*ServerReplica{
							0: {loadedModels: map[string]bool{}},
							1: {loadedModels: map[string]bool{}},
						},
					},
				},
			},
			modelKey:        "model",
			version:         "1",
			serverKey:       "server",
			replicaIdx:      0,
			state:           Unloaded,
			loaded:          false,
			availableMemory: 20,
			deleted:         true,
			err:             nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := NewMemoryStore(logger, test.store)
			err := ms.UpdateModelState(test.modelKey, test.version, test.serverKey, test.replicaIdx, &test.availableMemory, test.state)
			if test.err == nil {
				g.Expect(err).To(BeNil())
				if !test.deleted {
					g.Expect(test.store.models[test.modelKey].Latest().GetModelReplicaState(test.replicaIdx)).To(Equal(test.state))
					g.Expect(test.store.servers[test.serverKey].replicas[test.replicaIdx].loadedModels[test.modelKey]).To(Equal(test.loaded))
				} else {
					g.Expect(test.store.models[test.modelKey]).To(BeNil())
				}

			} else {
				g.Expect(err).ToNot(BeNil())
				g.Expect(errors.Is(err, test.err)).To(BeTrue())
			}
		})
	}
}