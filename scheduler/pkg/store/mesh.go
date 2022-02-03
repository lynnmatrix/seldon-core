package store

import (
	"time"

	pba "github.com/seldonio/seldon-core/scheduler/apis/mlops/agent"
	pb "github.com/seldonio/seldon-core/scheduler/apis/mlops/scheduler"
	"google.golang.org/protobuf/proto"
)

type LocalSchedulerStore struct {
	servers                map[string]*Server
	models                 map[string]*Model
	failedToScheduleModels map[string]bool
}

func NewLocalSchedulerStore() *LocalSchedulerStore {
	m := LocalSchedulerStore{}
	m.servers = make(map[string]*Server)
	m.models = make(map[string]*Model)
	m.failedToScheduleModels = make(map[string]bool)
	return &m
}

type Model struct {
	versions []*ModelVersion
	deleted  bool
}

type ModelVersion struct {
	modelDefn *pb.Model
	version   uint32
	server    string
	replicas  map[int]ReplicaStatus
	deleted   bool
	state     ModelStatus
}

type ModelStatus struct {
	State               ModelState
	Reason              string
	AvailableReplicas   uint32
	UnavailableReplicas uint32
	Timestamp           time.Time
}

type ReplicaStatus struct {
	State     ModelReplicaState
	Reason    string
	Timestamp time.Time
}

func NewDefaultModelVersion(model *pb.Model, version uint32) *ModelVersion {
	return &ModelVersion{
		version:   version,
		modelDefn: model,
		replicas:  make(map[int]ReplicaStatus),
		deleted:   false,
		state:     ModelStatus{State: ModelStateUnknown},
	}
}

func NewModelVersion(model *pb.Model, version uint32, server string, replicas map[int]ReplicaStatus, deleted bool, state ModelState) *ModelVersion {
	return &ModelVersion{
		version:   version,
		modelDefn: model,
		server:    server,
		replicas:  replicas,
		deleted:   deleted,
		state:     ModelStatus{State: state},
	}
}

type Server struct {
	name             string
	replicas         map[int]*ServerReplica
	shared           bool
	expectedReplicas int
	kubernetesMeta   *pb.KubernetesMeta
}

func (s *Server) SetExpectedReplicas(replicas int) {
	s.expectedReplicas = replicas
}

func (s *Server) SetKubernetesMeta(meta *pb.KubernetesMeta) {
	s.kubernetesMeta = meta
}

func NewServer(name string, shared bool) *Server {
	return &Server{
		name:             name,
		replicas:         make(map[int]*ServerReplica),
		shared:           shared,
		expectedReplicas: -1,
	}
}

type ServerReplica struct {
	inferenceSvc      string
	inferenceHttpPort int32
	inferenceGrpcPort int32
	replicaIdx        int
	server            *Server
	capabilities      []string
	memory            uint64
	availableMemory   uint64
	loadedModels      map[string]bool
	overCommit        bool
}

func NewServerReplica(inferenceSvc string,
	inferenceHttpPort int32,
	inferenceGrpcPort int32,
	replicaIdx int,
	server *Server,
	capabilities []string,
	memory uint64,
	availableMemory uint64,
	loadedModels map[string]bool,
	overCommit bool) *ServerReplica {
	return &ServerReplica{
		inferenceSvc:      inferenceSvc,
		inferenceHttpPort: inferenceHttpPort,
		inferenceGrpcPort: inferenceGrpcPort,
		replicaIdx:        replicaIdx,
		server:            server,
		capabilities:      capabilities,
		memory:            memory,
		availableMemory:   availableMemory,
		loadedModels:      loadedModels,
		overCommit:        overCommit,
	}
}

func NewServerReplicaFromConfig(server *Server, replicaIdx int, loadedModels map[string]bool, config *pba.ReplicaConfig, availableMemoryBytes uint64) *ServerReplica {
	return &ServerReplica{
		inferenceSvc:      config.GetInferenceSvc(),
		inferenceHttpPort: config.GetInferenceHttpPort(),
		inferenceGrpcPort: config.GetInferenceGrpcPort(),
		replicaIdx:        replicaIdx,
		server:            server,
		capabilities:      config.GetCapabilities(),
		memory:            config.GetMemoryBytes(),
		availableMemory:   availableMemoryBytes,
		loadedModels:      loadedModels,
		overCommit:        config.GetOverCommit(),
	}
}

type ModelState uint32

const (
	ModelStateUnknown ModelState = iota
	ModelProgressing
	ModelAvailable
	ModelFailed
	ModelTerminating
	ModelTerminated
	ModelTerminateFailed
	ScheduleFailed
)

func (m ModelState) String() string {
	return [...]string{"ModelStateUnknown", "ModelProgressing", "ModelAvailable", "ModelFailed", "ModelTerminating", "ModelTerminated", "ModelTerminateFailed", "ScheduleFailed"}[m]
}

type ModelReplicaState uint32

const (
	ModelReplicaStateUnknown ModelReplicaState = iota
	LoadRequested
	Loading
	Loaded
	LoadFailed
	UnloadRequested
	Unloading
	Unloaded
	UnloadFailed
	Available
	LoadedUnavailable
)

var replicaStates = []ModelReplicaState{
	ModelReplicaStateUnknown,
	LoadRequested,
	Loading,
	Loaded,
	LoadFailed,
	UnloadRequested,
	Unloading,
	Unloaded,
	UnloadFailed,
	Available,
	LoadedUnavailable,
}

func (m ModelReplicaState) NoProgressingEndpoint() bool {
	return (m == Unloaded || m == ModelReplicaStateUnknown || m == UnloadFailed || m == Unloading || m == UnloadRequested)
}

func (m ModelReplicaState) AlreadyLoadingOrLoaded() bool {
	return (m == Loading || m == Loaded || m == Available || m == LoadedUnavailable)
}

func (m ModelReplicaState) AlreadyUnloadingOrUnloaded() bool {
	return (m == Unloading || m == Unloaded)
}

func (m ModelReplicaState) Inactive() bool {
	return (m == Unloaded || m == UnloadFailed || m == ModelReplicaStateUnknown)
}

func (m ModelReplicaState) IsLoadingOrLoaded() bool {
	return (m == Loaded || m == LoadRequested || m == Loading || m == Available || m == LoadedUnavailable)
}

func (me ModelReplicaState) String() string {
	return [...]string{"ModelReplicaStateUnknown", "LoadRequested", "Loading", "Loaded", "LoadFailed", "UnloadRequested", "Unloading", "Unloaded", "UnloadFailed", "Available", "LoadedUnavailable"}[me]
}

func (m *Model) HasLatest() bool {
	return len(m.versions) > 0
}

func (m *Model) Latest() *ModelVersion {
	if len(m.versions) > 0 {
		return m.versions[len(m.versions)-1]
	} else {
		return nil
	}
}

func (m *Model) GetVersion(version uint32) *ModelVersion {
	for _, mv := range m.versions {
		if mv.GetVersion() == version {
			return mv
		}
	}
	return nil
}

func (m *Model) GetVersions() []uint32 {
	versions := make([]uint32, len(m.versions))
	for idx, v := range m.versions {
		versions[idx] = v.version
	}
	return versions
}

func (m *Model) getLastAvailableModelVersionIdx() int {
	lastAvailableIdx := -1
	for idx, mv := range m.versions {
		if mv.state.State == ModelAvailable {
			lastAvailableIdx = idx
		}
	}
	return lastAvailableIdx
}

func (m *Model) GetLastAvailableModelVersion() *ModelVersion {
	lastAvailableIdx := m.getLastAvailableModelVersionIdx()
	if lastAvailableIdx != -1 {
		return m.versions[lastAvailableIdx]
	}
	return nil
}

func (m *Model) Previous() *ModelVersion {
	if len(m.versions) > 1 {
		return m.versions[len(m.versions)-2]
	} else {
		return nil
	}
}

//TODO do we need to consider previous versions?
func (m *Model) Inactive() bool {
	return m.Latest().Inactive()
}

func (m *Model) IsDeleted() bool {
	return m.deleted
}

func (m *ModelVersion) GetVersion() uint32 {
	return m.version
}

func (m *ModelVersion) GetRequiredMemory() uint64 {
	return m.modelDefn.GetModelSpec().GetMemoryBytes()
}

func (m *ModelVersion) GetRequirements() []string {
	return m.modelDefn.GetModelSpec().GetRequirements()
}

func (m *ModelVersion) DesiredReplicas() int {
	return int(m.modelDefn.GetDeploymentSpec().GetReplicas())
}

func (m *ModelVersion) GetModel() *pb.Model {
	return proto.Clone(m.modelDefn).(*pb.Model)
}

func (m *ModelVersion) GetMeta() *pb.MetaData {
	return proto.Clone(m.modelDefn.GetMeta()).(*pb.MetaData)
}

func (m *ModelVersion) GetModelSpec() *pb.ModelSpec {
	return proto.Clone(m.modelDefn.GetModelSpec()).(*pb.ModelSpec)
}

func (m *ModelVersion) GetDeploymentSpec() *pb.DeploymentSpec {
	return proto.Clone(m.modelDefn.GetDeploymentSpec()).(*pb.DeploymentSpec)
}

func (m *ModelVersion) SetDeploymentSpec(spec *pb.DeploymentSpec) {
	m.modelDefn.DeploymentSpec = spec
}

func (m *ModelVersion) Server() string {
	return m.server
}

func (m *ModelVersion) ReplicaState() map[int]ReplicaStatus {
	return m.replicas
}

func (m *ModelVersion) ModelState() ModelStatus {
	return m.state
}

func (m *ModelVersion) GetModelReplicaState(replicaIdx int) ModelReplicaState {
	state, ok := m.replicas[replicaIdx]
	if !ok {
		return ModelReplicaStateUnknown
	}
	return state.State
}

func (m *ModelVersion) UpdateKubernetesMeta(meta *pb.KubernetesMeta) {
	m.modelDefn.Meta.KubernetesMeta = meta
}

func (m *ModelVersion) GetReplicaForState(state ModelReplicaState) []int {
	var assignment []int
	for k, v := range m.replicas {
		if v.State == state {
			assignment = append(assignment, k)
		}
	}
	return assignment
}

func (m *ModelVersion) GetRequestedServer() *string {
	return m.modelDefn.GetModelSpec().Server
}

func (m *ModelVersion) HasServer() bool {
	return m.server != ""
}

func (m *ModelVersion) Inactive() bool {
	for _, v := range m.replicas {
		if !v.State.Inactive() {
			return false
		}
	}
	return true
}

func (m *ModelVersion) IsLoadingOrLoaded(server string, replicaIdx int) bool {
	if server != m.server {
		return false
	}
	for r, v := range m.replicas {
		if r == replicaIdx && v.State.IsLoadingOrLoaded() {
			return true
		}
	}
	return false
}

func (m *ModelVersion) NoLiveReplica() bool {
	for _, v := range m.replicas {
		if !v.State.NoProgressingEndpoint() {
			return false
		}
	}
	return true
}

func (m *ModelVersion) GetAssignment() []int {
	var assignment []int
	for k, v := range m.replicas {
		if v.State == Loaded || v.State == Available || v.State == LoadedUnavailable {
			assignment = append(assignment, k)
		}
	}
	return assignment
}

func (m *ModelVersion) Key() string {
	return m.modelDefn.GetMeta().GetName()
}

func (m *ModelVersion) IsDeleted() bool {
	return m.deleted
}

func (m *ModelVersion) SetReplicaState(replicaIdx int, state ReplicaStatus) {
	m.replicas[replicaIdx] = state
}

func (s *Server) Key() string {
	return s.name
}

func (s *Server) NumReplicas() uint32 {
	return uint32(len(s.replicas))
}

func (s *Server) GetAvailableMemory(idx int) uint64 {
	if s != nil && idx < len(s.replicas) {
		return s.replicas[idx].availableMemory
	}
	return 0
}

func (s *Server) GetMemory(idx int) uint64 {
	if s != nil && idx < len(s.replicas) {
		return s.replicas[idx].memory
	}
	return 0
}

func (s *Server) GetReplicaInferenceSvc(idx int) string {
	return s.replicas[idx].inferenceSvc
}

func (s *Server) GetReplicaInferenceHttpPort(idx int) int32 {
	return s.replicas[idx].inferenceHttpPort
}

func (s *ServerReplica) GetLoadedModels() []string {
	var models []string
	for model := range s.loadedModels {
		models = append(models, model)
	}
	return models
}

func (s *ServerReplica) GetAvailableMemory() uint64 {
	return s.availableMemory
}

func (s *ServerReplica) GetMemory() uint64 {
	return s.memory
}

func (s *ServerReplica) GetCapabilities() []string {
	return s.capabilities
}

func (s *ServerReplica) GetReplicaIdx() int {
	return s.replicaIdx
}

func (s *ServerReplica) GetInferenceSvc() string {
	return s.inferenceSvc
}

func (s *ServerReplica) GetInferenceHttpPort() int32 {
	return s.inferenceHttpPort
}

func (s *ServerReplica) GetInferenceGrpcPort() int32 {
	return s.inferenceGrpcPort
}