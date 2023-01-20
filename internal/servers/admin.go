package servers

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/servers/mappers"
	"github.com/tupyy/tinyedge-controller/internal/services/configuration"
	"github.com/tupyy/tinyedge-controller/internal/services/device"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"github.com/tupyy/tinyedge-controller/internal/services/manifest"
	"github.com/tupyy/tinyedge-controller/internal/services/repository"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	pb "github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminServer struct {
	pb.UnsafeAdminServiceServer
	repositoryService *repository.Service
	manifestService   *manifest.Service
	deviceService     *device.Service
	confService       *configuration.Service
}

func NewAdminServer(r *repository.Service, m *manifest.Service, d *device.Service, c *configuration.Service) *AdminServer {
	return &AdminServer{repositoryService: r, manifestService: m, deviceService: d, confService: c}
}

func (a *AdminServer) GetDevices(ctx context.Context, req *pb.DevicesListRequest) (*pb.DevicesListResponse, error) {
	devices, err := a.deviceService.GetDevices(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	models := make([]*common.Device, 0, len(devices))
	for _, d := range devices {
		models = append(models, mappers.DeviceToProto(d))
	}

	return &pb.DevicesListResponse{
		Devices: models,
		Size:    int32(len(models)),
		Total:   int32(len(models)),
		Page:    1,
	}, nil
}

func (a *AdminServer) GetDevice(ctx context.Context, req *pb.IdRequest) (*common.Device, error) {
	device, err := a.deviceService.GetDevice(ctx, req.Id)
	if err != nil {
		if errService.IsResourceNotFound(err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return mappers.DeviceToProto(device), nil
}

func (a *AdminServer) UpdateDevice(ctx context.Context, req *pb.UpdateDeviceRequest) (*common.Device, error) {
	device, err := a.deviceService.GetDevice(ctx, req.DeviceId)
	if err != nil {
		if errService.IsResourceNotFound(err) {
			return &common.Device{}, status.Errorf(codes.NotFound, err.Error())
		}
		return &common.Device{}, status.Error(codes.Internal, "internal error")
	}

	if req.SetId != "" {
		_, err := a.deviceService.GetSet(ctx, req.SetId)
		if err != nil {
			if errService.IsResourceNotFound(err) {
				return &common.Device{}, status.Errorf(codes.NotFound, err.Error())
			}
			return &common.Device{}, status.Errorf(codes.Internal, "internal error")
		}
		device.SetID = &req.SetId
	}

	if req.NamespaceId != "" {
		_, err := a.deviceService.GetNamespace(ctx, req.NamespaceId)
		if err != nil {
			if errService.IsResourceNotFound(err) {
				return &common.Device{}, status.Errorf(codes.NotFound, err.Error())
			}
			return &common.Device{}, status.Errorf(codes.Internal, "internal error")
		}
		device.NamespaceID = req.NamespaceId
	}

	if req.ConfigurationId != "" {
		_, err := a.confService.GetConfiguration(ctx, req.ConfigurationId)
		if err != nil {
			if errService.IsResourceNotFound(err) {
				return &common.Device{}, status.Errorf(codes.NotFound, err.Error())
			}
			return &common.Device{}, status.Errorf(codes.Internal, "internal error")
		}
		device.Configuration = &entity.Configuration{
			ID: req.ConfigurationId,
		}
	}

	if err := a.deviceService.UpdateDevice(ctx, device); err != nil {
		return &common.Device{}, fmt.Errorf("unable to update device %q", device.ID)
	}

	return mappers.DeviceToProto(device), nil
}

func (a *AdminServer) AddSet(ctx context.Context, req *pb.AddSetRequest) (*common.Set, error) {
	if req.SetName == "" || req.NamespaceId == "" {
		return nil, status.Error(codes.InvalidArgument, "set name or namespace id is missing")
	}
	set := entity.Set{
		Name:        req.SetName,
		NamespaceID: req.NamespaceId,
	}
	if req.ConfigurationId != nil {
		set.Configuration = &entity.Configuration{
			ID: *req.ConfigurationId,
		}
	}

	err := a.deviceService.CreateSet(ctx, set)
	if errService.IsResourceAlreadyExists(err) {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	} else if errService.IsResourceNotFound(err) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	pbSet := &common.Set{Name: req.SetName, Namespace: req.NamespaceId}
	if req.ConfigurationId != nil {
		pbSet.Configuration = *req.ConfigurationId
	}

	return pbSet, nil
}

func (a *AdminServer) AddNamespace(ctx context.Context, req *pb.AddNamespaceRequest) (*pb.Namespace, error) {
	if req.Name == "" || req.ConfigurationId == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace name or configuration id is missing")
	}
	err := a.deviceService.CreateNamespace(ctx, entity.Namespace{
		Name: req.Name,
		Configuration: entity.Configuration{
			ID: req.ConfigurationId,
		},
		IsDefault: req.IsDefault,
	})

	if errService.IsResourceAlreadyExists(err) {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.Namespace{Name: req.Name, Configuration: req.ConfigurationId, IsDefault: req.IsDefault}, nil
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
	set, err := a.deviceService.GetSet(ctx, req.Id)
	if err != nil {
		if errService.IsResourceNotFound(err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return mappers.SetToProto(set), nil
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

	if len(models) == 0 {
		return &pb.NamespaceListResponse{}, nil
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

	if len(pgManifests) == 0 {
		return &pb.ManifestListResponse{}, nil
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
	manifest, err := a.manifestService.GetManifest(ctx, req.Id)
	if err != nil {
		if errService.IsResourceNotFound(err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return mappers.ManifestToProto(manifest), nil
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

	if len(models) == 0 {
		return &pb.RepositoryListResponse{}, nil
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
	repo := entity.Repository{
		Url: req.Url,
		Id:  req.Name,
	}

	if err := a.repositoryService.Add(ctx, repo); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unable to add repository %s", err)
	}

	return &pb.AddRepositoryResponse{
		Url:  req.Url,
		Name: req.Name,
	}, nil
}
