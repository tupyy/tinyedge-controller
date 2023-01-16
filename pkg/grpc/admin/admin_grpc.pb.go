// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: admin.proto

package admin

import (
	context "context"
	common "github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AdminServiceClient is the client API for AdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminServiceClient interface {
	// GetDevices returns a list of devices.
	GetDevices(ctx context.Context, in *DevicesListRequest, opts ...grpc.CallOption) (*DevicesListResponse, error)
	// GetDevice returns a device.
	GetDevice(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*common.Device, error)
	// AddWorkloadToSet add a device to a set.
	AddDeviceToSet(ctx context.Context, in *DeviceToSetRequest, opts ...grpc.CallOption) (*common.Empty, error)
	// RemoveDeviceFromSet removes a device from a set.
	RemoveDeviceFromSet(ctx context.Context, in *DeviceToSetRequest, opts ...grpc.CallOption) (*common.Empty, error)
	// GetSets returns a list of device sets.
	GetSets(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*SetsListResponse, error)
	// GetSet returns a device set.
	GetSet(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*common.Set, error)
	// AddSet adds a set
	AddSet(ctx context.Context, in *AddSetRequest, opts ...grpc.CallOption) (*common.Set, error)
	// GetNamespaces returns a list with namespaces
	GetNamespaces(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*NamespaceListResponse, error)
	// GetManifests return a list of manifests
	GetManifests(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ManifestListResponse, error)
	// GetManifest return a manifests
	GetManifest(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Manifest, error)
	// GetRepositories return a list of repositories
	GetRepositories(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*RepositoryListResponse, error)
	// AddRepository add a repository
	AddRepository(ctx context.Context, in *AddRepositoryRequest, opts ...grpc.CallOption) (*AddRepositoryResponse, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) GetDevices(ctx context.Context, in *DevicesListRequest, opts ...grpc.CallOption) (*DevicesListResponse, error) {
	out := new(DevicesListResponse)
	err := c.cc.Invoke(ctx, "/AdminService/GetDevices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetDevice(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*common.Device, error) {
	out := new(common.Device)
	err := c.cc.Invoke(ctx, "/AdminService/GetDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) AddDeviceToSet(ctx context.Context, in *DeviceToSetRequest, opts ...grpc.CallOption) (*common.Empty, error) {
	out := new(common.Empty)
	err := c.cc.Invoke(ctx, "/AdminService/AddDeviceToSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) RemoveDeviceFromSet(ctx context.Context, in *DeviceToSetRequest, opts ...grpc.CallOption) (*common.Empty, error) {
	out := new(common.Empty)
	err := c.cc.Invoke(ctx, "/AdminService/RemoveDeviceFromSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetSets(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*SetsListResponse, error) {
	out := new(SetsListResponse)
	err := c.cc.Invoke(ctx, "/AdminService/GetSets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetSet(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*common.Set, error) {
	out := new(common.Set)
	err := c.cc.Invoke(ctx, "/AdminService/GetSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) AddSet(ctx context.Context, in *AddSetRequest, opts ...grpc.CallOption) (*common.Set, error) {
	out := new(common.Set)
	err := c.cc.Invoke(ctx, "/AdminService/AddSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetNamespaces(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*NamespaceListResponse, error) {
	out := new(NamespaceListResponse)
	err := c.cc.Invoke(ctx, "/AdminService/GetNamespaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetManifests(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ManifestListResponse, error) {
	out := new(ManifestListResponse)
	err := c.cc.Invoke(ctx, "/AdminService/GetManifests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetManifest(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Manifest, error) {
	out := new(Manifest)
	err := c.cc.Invoke(ctx, "/AdminService/GetManifest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetRepositories(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*RepositoryListResponse, error) {
	out := new(RepositoryListResponse)
	err := c.cc.Invoke(ctx, "/AdminService/GetRepositories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) AddRepository(ctx context.Context, in *AddRepositoryRequest, opts ...grpc.CallOption) (*AddRepositoryResponse, error) {
	out := new(AddRepositoryResponse)
	err := c.cc.Invoke(ctx, "/AdminService/AddRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServiceServer is the server API for AdminService service.
// All implementations must embed UnimplementedAdminServiceServer
// for forward compatibility
type AdminServiceServer interface {
	// GetDevices returns a list of devices.
	GetDevices(context.Context, *DevicesListRequest) (*DevicesListResponse, error)
	// GetDevice returns a device.
	GetDevice(context.Context, *IdRequest) (*common.Device, error)
	// AddWorkloadToSet add a device to a set.
	AddDeviceToSet(context.Context, *DeviceToSetRequest) (*common.Empty, error)
	// RemoveDeviceFromSet removes a device from a set.
	RemoveDeviceFromSet(context.Context, *DeviceToSetRequest) (*common.Empty, error)
	// GetSets returns a list of device sets.
	GetSets(context.Context, *ListRequest) (*SetsListResponse, error)
	// GetSet returns a device set.
	GetSet(context.Context, *IdRequest) (*common.Set, error)
	// AddSet adds a set
	AddSet(context.Context, *AddSetRequest) (*common.Set, error)
	// GetNamespaces returns a list with namespaces
	GetNamespaces(context.Context, *ListRequest) (*NamespaceListResponse, error)
	// GetManifests return a list of manifests
	GetManifests(context.Context, *ListRequest) (*ManifestListResponse, error)
	// GetManifest return a manifests
	GetManifest(context.Context, *IdRequest) (*Manifest, error)
	// GetRepositories return a list of repositories
	GetRepositories(context.Context, *ListRequest) (*RepositoryListResponse, error)
	// AddRepository add a repository
	AddRepository(context.Context, *AddRepositoryRequest) (*AddRepositoryResponse, error)
	mustEmbedUnimplementedAdminServiceServer()
}

// UnimplementedAdminServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAdminServiceServer struct {
}

func (UnimplementedAdminServiceServer) GetDevices(context.Context, *DevicesListRequest) (*DevicesListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDevices not implemented")
}
func (UnimplementedAdminServiceServer) GetDevice(context.Context, *IdRequest) (*common.Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDevice not implemented")
}
func (UnimplementedAdminServiceServer) AddDeviceToSet(context.Context, *DeviceToSetRequest) (*common.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDeviceToSet not implemented")
}
func (UnimplementedAdminServiceServer) RemoveDeviceFromSet(context.Context, *DeviceToSetRequest) (*common.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveDeviceFromSet not implemented")
}
func (UnimplementedAdminServiceServer) GetSets(context.Context, *ListRequest) (*SetsListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSets not implemented")
}
func (UnimplementedAdminServiceServer) GetSet(context.Context, *IdRequest) (*common.Set, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSet not implemented")
}
func (UnimplementedAdminServiceServer) AddSet(context.Context, *AddSetRequest) (*common.Set, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSet not implemented")
}
func (UnimplementedAdminServiceServer) GetNamespaces(context.Context, *ListRequest) (*NamespaceListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNamespaces not implemented")
}
func (UnimplementedAdminServiceServer) GetManifests(context.Context, *ListRequest) (*ManifestListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetManifests not implemented")
}
func (UnimplementedAdminServiceServer) GetManifest(context.Context, *IdRequest) (*Manifest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetManifest not implemented")
}
func (UnimplementedAdminServiceServer) GetRepositories(context.Context, *ListRequest) (*RepositoryListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRepositories not implemented")
}
func (UnimplementedAdminServiceServer) AddRepository(context.Context, *AddRepositoryRequest) (*AddRepositoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRepository not implemented")
}
func (UnimplementedAdminServiceServer) mustEmbedUnimplementedAdminServiceServer() {}

// UnsafeAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServiceServer will
// result in compilation errors.
type UnsafeAdminServiceServer interface {
	mustEmbedUnimplementedAdminServiceServer()
}

func RegisterAdminServiceServer(s grpc.ServiceRegistrar, srv AdminServiceServer) {
	s.RegisterService(&AdminService_ServiceDesc, srv)
}

func _AdminService_GetDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DevicesListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetDevices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetDevices(ctx, req.(*DevicesListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetDevice(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_AddDeviceToSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeviceToSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).AddDeviceToSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/AddDeviceToSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).AddDeviceToSet(ctx, req.(*DeviceToSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_RemoveDeviceFromSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeviceToSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).RemoveDeviceFromSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/RemoveDeviceFromSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).RemoveDeviceFromSet(ctx, req.(*DeviceToSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetSets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetSets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetSets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetSets(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetSet(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_AddSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).AddSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/AddSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).AddSet(ctx, req.(*AddSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetNamespaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetNamespaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetNamespaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetNamespaces(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetManifests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetManifests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetManifests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetManifests(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetManifest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetManifest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetManifest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetManifest(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetRepositories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetRepositories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetRepositories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetRepositories(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_AddRepository_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRepositoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).AddRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/AddRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).AddRepository(ctx, req.(*AddRepositoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AdminService_ServiceDesc is the grpc.ServiceDesc for AdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDevices",
			Handler:    _AdminService_GetDevices_Handler,
		},
		{
			MethodName: "GetDevice",
			Handler:    _AdminService_GetDevice_Handler,
		},
		{
			MethodName: "AddDeviceToSet",
			Handler:    _AdminService_AddDeviceToSet_Handler,
		},
		{
			MethodName: "RemoveDeviceFromSet",
			Handler:    _AdminService_RemoveDeviceFromSet_Handler,
		},
		{
			MethodName: "GetSets",
			Handler:    _AdminService_GetSets_Handler,
		},
		{
			MethodName: "GetSet",
			Handler:    _AdminService_GetSet_Handler,
		},
		{
			MethodName: "AddSet",
			Handler:    _AdminService_AddSet_Handler,
		},
		{
			MethodName: "GetNamespaces",
			Handler:    _AdminService_GetNamespaces_Handler,
		},
		{
			MethodName: "GetManifests",
			Handler:    _AdminService_GetManifests_Handler,
		},
		{
			MethodName: "GetManifest",
			Handler:    _AdminService_GetManifest_Handler,
		},
		{
			MethodName: "GetRepositories",
			Handler:    _AdminService_GetRepositories_Handler,
		},
		{
			MethodName: "AddRepository",
			Handler:    _AdminService_AddRepository_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "admin.proto",
}
