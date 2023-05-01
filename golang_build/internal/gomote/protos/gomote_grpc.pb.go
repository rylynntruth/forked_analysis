// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: gomote.proto

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GomoteServiceClient is the client API for GomoteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GomoteServiceClient interface {
	// Authenticate provides authentication information without any additional action.
	Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error)
	// AddBootstrap adds the bootstrap version of Go to the work directory.
	AddBootstrap(ctx context.Context, in *AddBootstrapRequest, opts ...grpc.CallOption) (*AddBootstrapResponse, error)
	// CreateInstance creates a gomote instance.
	CreateInstance(ctx context.Context, in *CreateInstanceRequest, opts ...grpc.CallOption) (GomoteService_CreateInstanceClient, error)
	// DestroyInstance destroys a gomote instance.
	DestroyInstance(ctx context.Context, in *DestroyInstanceRequest, opts ...grpc.CallOption) (*DestroyInstanceResponse, error)
	// ExecuteCommand executes a command on the gomote instance.
	ExecuteCommand(ctx context.Context, in *ExecuteCommandRequest, opts ...grpc.CallOption) (GomoteService_ExecuteCommandClient, error)
	// InstanceAlive gives the liveness state of a gomote instance.
	InstanceAlive(ctx context.Context, in *InstanceAliveRequest, opts ...grpc.CallOption) (*InstanceAliveResponse, error)
	// ListDirectory lists the contents of a directory on an gomote instance.
	ListDirectory(ctx context.Context, in *ListDirectoryRequest, opts ...grpc.CallOption) (*ListDirectoryResponse, error)
	// ListInstances lists all of the live gomote instances owned by the caller.
	ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesResponse, error)
	// ReadTGZToURL tars and zips a directory which exists on the gomote instance and returns a URL where it can be
	// downloaded from.
	ReadTGZToURL(ctx context.Context, in *ReadTGZToURLRequest, opts ...grpc.CallOption) (*ReadTGZToURLResponse, error)
	// RemoveFiles removes files or directories from the gomote instance.
	RemoveFiles(ctx context.Context, in *RemoveFilesRequest, opts ...grpc.CallOption) (*RemoveFilesResponse, error)
	// SignSSHKey signs an SSH public key which can be used to SSH into instances owned by the caller.
	SignSSHKey(ctx context.Context, in *SignSSHKeyRequest, opts ...grpc.CallOption) (*SignSSHKeyResponse, error)
	// UploadFile generates a signed URL and associated fields to be used when uploading the object to GCS. Once uploaded
	// the corresponding Write endpoint can be used to send the file to the gomote instance.
	UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error)
	// WriteFileFromURL
	WriteFileFromURL(ctx context.Context, in *WriteFileFromURLRequest, opts ...grpc.CallOption) (*WriteFileFromURLResponse, error)
	// WriteTGZFromURL retrieves a tar and zipped file from a URL and expands it onto the file system of a gomote instance.
	WriteTGZFromURL(ctx context.Context, in *WriteTGZFromURLRequest, opts ...grpc.CallOption) (*WriteTGZFromURLResponse, error)
}

type gomoteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGomoteServiceClient(cc grpc.ClientConnInterface) GomoteServiceClient {
	return &gomoteServiceClient{cc}
}

func (c *gomoteServiceClient) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/Authenticate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) AddBootstrap(ctx context.Context, in *AddBootstrapRequest, opts ...grpc.CallOption) (*AddBootstrapResponse, error) {
	out := new(AddBootstrapResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/AddBootstrap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) CreateInstance(ctx context.Context, in *CreateInstanceRequest, opts ...grpc.CallOption) (GomoteService_CreateInstanceClient, error) {
	stream, err := c.cc.NewStream(ctx, &GomoteService_ServiceDesc.Streams[0], "/protos.GomoteService/CreateInstance", opts...)
	if err != nil {
		return nil, err
	}
	x := &gomoteServiceCreateInstanceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GomoteService_CreateInstanceClient interface {
	Recv() (*CreateInstanceResponse, error)
	grpc.ClientStream
}

type gomoteServiceCreateInstanceClient struct {
	grpc.ClientStream
}

func (x *gomoteServiceCreateInstanceClient) Recv() (*CreateInstanceResponse, error) {
	m := new(CreateInstanceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gomoteServiceClient) DestroyInstance(ctx context.Context, in *DestroyInstanceRequest, opts ...grpc.CallOption) (*DestroyInstanceResponse, error) {
	out := new(DestroyInstanceResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/DestroyInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) ExecuteCommand(ctx context.Context, in *ExecuteCommandRequest, opts ...grpc.CallOption) (GomoteService_ExecuteCommandClient, error) {
	stream, err := c.cc.NewStream(ctx, &GomoteService_ServiceDesc.Streams[1], "/protos.GomoteService/ExecuteCommand", opts...)
	if err != nil {
		return nil, err
	}
	x := &gomoteServiceExecuteCommandClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GomoteService_ExecuteCommandClient interface {
	Recv() (*ExecuteCommandResponse, error)
	grpc.ClientStream
}

type gomoteServiceExecuteCommandClient struct {
	grpc.ClientStream
}

func (x *gomoteServiceExecuteCommandClient) Recv() (*ExecuteCommandResponse, error) {
	m := new(ExecuteCommandResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gomoteServiceClient) InstanceAlive(ctx context.Context, in *InstanceAliveRequest, opts ...grpc.CallOption) (*InstanceAliveResponse, error) {
	out := new(InstanceAliveResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/InstanceAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) ListDirectory(ctx context.Context, in *ListDirectoryRequest, opts ...grpc.CallOption) (*ListDirectoryResponse, error) {
	out := new(ListDirectoryResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/ListDirectory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesResponse, error) {
	out := new(ListInstancesResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/ListInstances", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) ReadTGZToURL(ctx context.Context, in *ReadTGZToURLRequest, opts ...grpc.CallOption) (*ReadTGZToURLResponse, error) {
	out := new(ReadTGZToURLResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/ReadTGZToURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) RemoveFiles(ctx context.Context, in *RemoveFilesRequest, opts ...grpc.CallOption) (*RemoveFilesResponse, error) {
	out := new(RemoveFilesResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/RemoveFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) SignSSHKey(ctx context.Context, in *SignSSHKeyRequest, opts ...grpc.CallOption) (*SignSSHKeyResponse, error) {
	out := new(SignSSHKeyResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/SignSSHKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error) {
	out := new(UploadFileResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/UploadFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) WriteFileFromURL(ctx context.Context, in *WriteFileFromURLRequest, opts ...grpc.CallOption) (*WriteFileFromURLResponse, error) {
	out := new(WriteFileFromURLResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/WriteFileFromURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gomoteServiceClient) WriteTGZFromURL(ctx context.Context, in *WriteTGZFromURLRequest, opts ...grpc.CallOption) (*WriteTGZFromURLResponse, error) {
	out := new(WriteTGZFromURLResponse)
	err := c.cc.Invoke(ctx, "/protos.GomoteService/WriteTGZFromURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GomoteServiceServer is the server API for GomoteService service.
// All implementations must embed UnimplementedGomoteServiceServer
// for forward compatibility
type GomoteServiceServer interface {
	// Authenticate provides authentication information without any additional action.
	Authenticate(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error)
	// AddBootstrap adds the bootstrap version of Go to the work directory.
	AddBootstrap(context.Context, *AddBootstrapRequest) (*AddBootstrapResponse, error)
	// CreateInstance creates a gomote instance.
	CreateInstance(*CreateInstanceRequest, GomoteService_CreateInstanceServer) error
	// DestroyInstance destroys a gomote instance.
	DestroyInstance(context.Context, *DestroyInstanceRequest) (*DestroyInstanceResponse, error)
	// ExecuteCommand executes a command on the gomote instance.
	ExecuteCommand(*ExecuteCommandRequest, GomoteService_ExecuteCommandServer) error
	// InstanceAlive gives the liveness state of a gomote instance.
	InstanceAlive(context.Context, *InstanceAliveRequest) (*InstanceAliveResponse, error)
	// ListDirectory lists the contents of a directory on an gomote instance.
	ListDirectory(context.Context, *ListDirectoryRequest) (*ListDirectoryResponse, error)
	// ListInstances lists all of the live gomote instances owned by the caller.
	ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesResponse, error)
	// ReadTGZToURL tars and zips a directory which exists on the gomote instance and returns a URL where it can be
	// downloaded from.
	ReadTGZToURL(context.Context, *ReadTGZToURLRequest) (*ReadTGZToURLResponse, error)
	// RemoveFiles removes files or directories from the gomote instance.
	RemoveFiles(context.Context, *RemoveFilesRequest) (*RemoveFilesResponse, error)
	// SignSSHKey signs an SSH public key which can be used to SSH into instances owned by the caller.
	SignSSHKey(context.Context, *SignSSHKeyRequest) (*SignSSHKeyResponse, error)
	// UploadFile generates a signed URL and associated fields to be used when uploading the object to GCS. Once uploaded
	// the corresponding Write endpoint can be used to send the file to the gomote instance.
	UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error)
	// WriteFileFromURL
	WriteFileFromURL(context.Context, *WriteFileFromURLRequest) (*WriteFileFromURLResponse, error)
	// WriteTGZFromURL retrieves a tar and zipped file from a URL and expands it onto the file system of a gomote instance.
	WriteTGZFromURL(context.Context, *WriteTGZFromURLRequest) (*WriteTGZFromURLResponse, error)
	mustEmbedUnimplementedGomoteServiceServer()
}

// UnimplementedGomoteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGomoteServiceServer struct {
}

func (UnimplementedGomoteServiceServer) Authenticate(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedGomoteServiceServer) AddBootstrap(context.Context, *AddBootstrapRequest) (*AddBootstrapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBootstrap not implemented")
}
func (UnimplementedGomoteServiceServer) CreateInstance(*CreateInstanceRequest, GomoteService_CreateInstanceServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateInstance not implemented")
}
func (UnimplementedGomoteServiceServer) DestroyInstance(context.Context, *DestroyInstanceRequest) (*DestroyInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DestroyInstance not implemented")
}
func (UnimplementedGomoteServiceServer) ExecuteCommand(*ExecuteCommandRequest, GomoteService_ExecuteCommandServer) error {
	return status.Errorf(codes.Unimplemented, "method ExecuteCommand not implemented")
}
func (UnimplementedGomoteServiceServer) InstanceAlive(context.Context, *InstanceAliveRequest) (*InstanceAliveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstanceAlive not implemented")
}
func (UnimplementedGomoteServiceServer) ListDirectory(context.Context, *ListDirectoryRequest) (*ListDirectoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDirectory not implemented")
}
func (UnimplementedGomoteServiceServer) ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInstances not implemented")
}
func (UnimplementedGomoteServiceServer) ReadTGZToURL(context.Context, *ReadTGZToURLRequest) (*ReadTGZToURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadTGZToURL not implemented")
}
func (UnimplementedGomoteServiceServer) RemoveFiles(context.Context, *RemoveFilesRequest) (*RemoveFilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFiles not implemented")
}
func (UnimplementedGomoteServiceServer) SignSSHKey(context.Context, *SignSSHKeyRequest) (*SignSSHKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignSSHKey not implemented")
}
func (UnimplementedGomoteServiceServer) UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedGomoteServiceServer) WriteFileFromURL(context.Context, *WriteFileFromURLRequest) (*WriteFileFromURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteFileFromURL not implemented")
}
func (UnimplementedGomoteServiceServer) WriteTGZFromURL(context.Context, *WriteTGZFromURLRequest) (*WriteTGZFromURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteTGZFromURL not implemented")
}
func (UnimplementedGomoteServiceServer) mustEmbedUnimplementedGomoteServiceServer() {}

// UnsafeGomoteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GomoteServiceServer will
// result in compilation errors.
type UnsafeGomoteServiceServer interface {
	mustEmbedUnimplementedGomoteServiceServer()
}

func RegisterGomoteServiceServer(s grpc.ServiceRegistrar, srv GomoteServiceServer) {
	s.RegisterService(&GomoteService_ServiceDesc, srv)
}

func _GomoteService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).Authenticate(ctx, req.(*AuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_AddBootstrap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBootstrapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).AddBootstrap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/AddBootstrap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).AddBootstrap(ctx, req.(*AddBootstrapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_CreateInstance_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CreateInstanceRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GomoteServiceServer).CreateInstance(m, &gomoteServiceCreateInstanceServer{stream})
}

type GomoteService_CreateInstanceServer interface {
	Send(*CreateInstanceResponse) error
	grpc.ServerStream
}

type gomoteServiceCreateInstanceServer struct {
	grpc.ServerStream
}

func (x *gomoteServiceCreateInstanceServer) Send(m *CreateInstanceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GomoteService_DestroyInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DestroyInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).DestroyInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/DestroyInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).DestroyInstance(ctx, req.(*DestroyInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_ExecuteCommand_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExecuteCommandRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GomoteServiceServer).ExecuteCommand(m, &gomoteServiceExecuteCommandServer{stream})
}

type GomoteService_ExecuteCommandServer interface {
	Send(*ExecuteCommandResponse) error
	grpc.ServerStream
}

type gomoteServiceExecuteCommandServer struct {
	grpc.ServerStream
}

func (x *gomoteServiceExecuteCommandServer) Send(m *ExecuteCommandResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GomoteService_InstanceAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstanceAliveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).InstanceAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/InstanceAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).InstanceAlive(ctx, req.(*InstanceAliveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_ListDirectory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDirectoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).ListDirectory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/ListDirectory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).ListDirectory(ctx, req.(*ListDirectoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_ListInstances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInstancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).ListInstances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/ListInstances",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).ListInstances(ctx, req.(*ListInstancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_ReadTGZToURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadTGZToURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).ReadTGZToURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/ReadTGZToURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).ReadTGZToURL(ctx, req.(*ReadTGZToURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_RemoveFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).RemoveFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/RemoveFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).RemoveFiles(ctx, req.(*RemoveFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_SignSSHKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignSSHKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).SignSSHKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/SignSSHKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).SignSSHKey(ctx, req.(*SignSSHKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_UploadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).UploadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/UploadFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).UploadFile(ctx, req.(*UploadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_WriteFileFromURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteFileFromURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).WriteFileFromURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/WriteFileFromURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).WriteFileFromURL(ctx, req.(*WriteFileFromURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GomoteService_WriteTGZFromURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteTGZFromURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GomoteServiceServer).WriteTGZFromURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.GomoteService/WriteTGZFromURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GomoteServiceServer).WriteTGZFromURL(ctx, req.(*WriteTGZFromURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GomoteService_ServiceDesc is the grpc.ServiceDesc for GomoteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GomoteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.GomoteService",
	HandlerType: (*GomoteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticate",
			Handler:    _GomoteService_Authenticate_Handler,
		},
		{
			MethodName: "AddBootstrap",
			Handler:    _GomoteService_AddBootstrap_Handler,
		},
		{
			MethodName: "DestroyInstance",
			Handler:    _GomoteService_DestroyInstance_Handler,
		},
		{
			MethodName: "InstanceAlive",
			Handler:    _GomoteService_InstanceAlive_Handler,
		},
		{
			MethodName: "ListDirectory",
			Handler:    _GomoteService_ListDirectory_Handler,
		},
		{
			MethodName: "ListInstances",
			Handler:    _GomoteService_ListInstances_Handler,
		},
		{
			MethodName: "ReadTGZToURL",
			Handler:    _GomoteService_ReadTGZToURL_Handler,
		},
		{
			MethodName: "RemoveFiles",
			Handler:    _GomoteService_RemoveFiles_Handler,
		},
		{
			MethodName: "SignSSHKey",
			Handler:    _GomoteService_SignSSHKey_Handler,
		},
		{
			MethodName: "UploadFile",
			Handler:    _GomoteService_UploadFile_Handler,
		},
		{
			MethodName: "WriteFileFromURL",
			Handler:    _GomoteService_WriteFileFromURL_Handler,
		},
		{
			MethodName: "WriteTGZFromURL",
			Handler:    _GomoteService_WriteTGZFromURL_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateInstance",
			Handler:       _GomoteService_CreateInstance_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ExecuteCommand",
			Handler:       _GomoteService_ExecuteCommand_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "gomote.proto",
}