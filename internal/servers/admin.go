package servers

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/servers/mappers"
	"github.com/tupyy/tinyedge-controller/internal/services/device"
	"github.com/tupyy/tinyedge-controller/internal/services/manifest"
	"github.com/tupyy/tinyedge-controller/internal/services/repository"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	pb "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	"go.uber.org/zap"
)

type AdminServer struct {
	pb.UnsafeAdminServiceServer
	repositoryService *repository.RepositoryService
	manifestService   *manifest.Service
	deviceService     *device.Service
}

func NewAdminServer(r *repository.RepositoryService, m *manifest.Service, d *device.Service) *AdminServer {
	return &AdminServer{repositoryService: r, manifestService: m, deviceService: d}
}

func (a *AdminServer) GetDevices(ctx context.Context, req *pb.DevicesListRequest) (*pb.DevicesListResponse, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

func (a *AdminServer) GetDevice(ctx context.Context, req *pb.IdRequest) (*common.Device, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

func (a *AdminServer) AddDeviceToSet(ctx context.Context, req *pb.DeviceToSetRequest) (*common.Empty, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

// RemoveDeviceFromSet removes a device from a set.
func (a *AdminServer) RemoveDeviceFromSet(ctx context.Context, req *pb.DeviceToSetRequest) (*common.Empty, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

func (a *AdminServer) AddSet(ctx context.Context, req *pb.AddSetRequest) (*common.Set, error) {
	return nil, nil
}

// GetDeviceSets returns a list of device sets.
func (a *AdminServer) GetSets(ctx context.Context, req *pb.ListRequest) (*pb.SetsListResponse, error) {
	sets, err := a.deviceService.GetSets(ctx)
	if err != nil {
		return nil, err
	}

	models := make([]*common.Set, 0, len(sets))
	for _, s := range sets {
		models = append(models, mappers.SetToProto(s))
	}

	return &pb.SetsListResponse{
		Sets:  models,
		Size:  int32(len(models)),
		Total: int32(len(models)),
		Page:  1,
	}, nil
}

// GetDeviceSet returns a device set.
func (a *AdminServer) GetSet(ctx context.Context, req *pb.IdRequest) (*common.Set, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

func (a *AdminServer) GetNamespaces(ctx context.Context, req *pb.ListRequest) (*pb.NamespaceListResponse, error) {
	namespaces, err := a.deviceService.GetNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	models := make([]*admin.Namespace, 0, len(namespaces))
	for _, n := range namespaces {
		models = append(models, mappers.NamespaceToProto(n))
	}

	return &pb.NamespaceListResponse{
		Namespaces: models,
		Size:       int32(len(models)),
		Total:      int32(len(models)),
		Page:       1,
	}, nil
}

// GetWorkloads return a list of workloads
func (a *AdminServer) GetManifests(ctx context.Context, req *pb.ListRequest) (*pb.ManifestListResponse, error) {
	repos, err := a.repositoryService.GetRepositories(ctx)
	if err != nil {
		return nil, err
	}
	manifests := make([]entity.ManifestWork, 0, 20)
	for _, r := range repos {
		m, err := a.manifestService.GetManifests(ctx, r)
		if err != nil {
			zap.S().Errorw("unable to get manifests from repository", "error", err, "repo_id", r.Id, "repo_url", r.Url)
			continue
		}
		manifests = append(manifests, m...)
	}

	pgManifests := make([]*pb.Manifest, 0, len(manifests))
	for _, m := range manifests {
		pgManifests = append(pgManifests, mappers.ManifestToProto(m))
	}

	return &pb.ManifestListResponse{
		Manifests: pgManifests,
		Size:      int32(len(pgManifests)),
		Total:     int32(len(pgManifests)),
		Page:      1,
	}, nil
}

// GetWorkload return a workload
func (a *AdminServer) GetManifest(ctx context.Context, req *pb.IdRequest) (*pb.Manifest, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

// GetRepositories return a list of repositories
func (a *AdminServer) GetRepositories(ctx context.Context, req *pb.ListRequest) (*pb.RepositoryListResponse, error) {
	repos, err := a.repositoryService.GetRepositories(ctx)
	if err != nil {
		return nil, err
	}

	models := make([]*admin.Repository, 0, len(repos))
	for _, r := range repos {
		models = append(models, mappers.RepositoryToModel(r))
	}

	return &pb.RepositoryListResponse{
		Repositories: models,
		Size:         int32(len(models)),
		Page:         1,
		Total:        int32(len(models)),
	}, nil
}

// AddRepository add a repository
func (a *AdminServer) AddRepository(ctx context.Context, req *pb.AddRepositoryRequest) (*pb.AddRepositoryResponse, error) {
	return nil, fmt.Errorf("unimplemented yet")
}
