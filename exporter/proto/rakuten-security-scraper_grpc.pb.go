// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: rakuten-security-scraper.proto

package proto

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

const (
	RakutenSecurityScraper_ListWithdrawalHistories_FullMethodName = "/RakutenSecurityScraper/ListWithdrawalHistories"
	RakutenSecurityScraper_ListDividendHistories_FullMethodName   = "/RakutenSecurityScraper/ListDividendHistories"
	RakutenSecurityScraper_TotalAssets_FullMethodName             = "/RakutenSecurityScraper/TotalAssets"
)

// RakutenSecurityScraperClient is the client API for RakutenSecurityScraper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RakutenSecurityScraperClient interface {
	ListWithdrawalHistories(ctx context.Context, in *ListWithdrawalHistoriesRequest, opts ...grpc.CallOption) (*ListWithdrawalHistoriesResponse, error)
	ListDividendHistories(ctx context.Context, in *ListDividendHistoriesRequest, opts ...grpc.CallOption) (*ListDividendHistoriesResponse, error)
	TotalAssets(ctx context.Context, in *TotalAssetRequest, opts ...grpc.CallOption) (*TotalAssetResponse, error)
}

type rakutenSecurityScraperClient struct {
	cc grpc.ClientConnInterface
}

func NewRakutenSecurityScraperClient(cc grpc.ClientConnInterface) RakutenSecurityScraperClient {
	return &rakutenSecurityScraperClient{cc}
}

func (c *rakutenSecurityScraperClient) ListWithdrawalHistories(ctx context.Context, in *ListWithdrawalHistoriesRequest, opts ...grpc.CallOption) (*ListWithdrawalHistoriesResponse, error) {
	out := new(ListWithdrawalHistoriesResponse)
	err := c.cc.Invoke(ctx, RakutenSecurityScraper_ListWithdrawalHistories_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rakutenSecurityScraperClient) ListDividendHistories(ctx context.Context, in *ListDividendHistoriesRequest, opts ...grpc.CallOption) (*ListDividendHistoriesResponse, error) {
	out := new(ListDividendHistoriesResponse)
	err := c.cc.Invoke(ctx, RakutenSecurityScraper_ListDividendHistories_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rakutenSecurityScraperClient) TotalAssets(ctx context.Context, in *TotalAssetRequest, opts ...grpc.CallOption) (*TotalAssetResponse, error) {
	out := new(TotalAssetResponse)
	err := c.cc.Invoke(ctx, RakutenSecurityScraper_TotalAssets_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RakutenSecurityScraperServer is the server API for RakutenSecurityScraper service.
// All implementations must embed UnimplementedRakutenSecurityScraperServer
// for forward compatibility
type RakutenSecurityScraperServer interface {
	ListWithdrawalHistories(context.Context, *ListWithdrawalHistoriesRequest) (*ListWithdrawalHistoriesResponse, error)
	ListDividendHistories(context.Context, *ListDividendHistoriesRequest) (*ListDividendHistoriesResponse, error)
	TotalAssets(context.Context, *TotalAssetRequest) (*TotalAssetResponse, error)
	mustEmbedUnimplementedRakutenSecurityScraperServer()
}

// UnimplementedRakutenSecurityScraperServer must be embedded to have forward compatible implementations.
type UnimplementedRakutenSecurityScraperServer struct {
}

func (UnimplementedRakutenSecurityScraperServer) ListWithdrawalHistories(context.Context, *ListWithdrawalHistoriesRequest) (*ListWithdrawalHistoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWithdrawalHistories not implemented")
}
func (UnimplementedRakutenSecurityScraperServer) ListDividendHistories(context.Context, *ListDividendHistoriesRequest) (*ListDividendHistoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDividendHistories not implemented")
}
func (UnimplementedRakutenSecurityScraperServer) TotalAssets(context.Context, *TotalAssetRequest) (*TotalAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TotalAssets not implemented")
}
func (UnimplementedRakutenSecurityScraperServer) mustEmbedUnimplementedRakutenSecurityScraperServer() {
}

// UnsafeRakutenSecurityScraperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RakutenSecurityScraperServer will
// result in compilation errors.
type UnsafeRakutenSecurityScraperServer interface {
	mustEmbedUnimplementedRakutenSecurityScraperServer()
}

func RegisterRakutenSecurityScraperServer(s grpc.ServiceRegistrar, srv RakutenSecurityScraperServer) {
	s.RegisterService(&RakutenSecurityScraper_ServiceDesc, srv)
}

func _RakutenSecurityScraper_ListWithdrawalHistories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWithdrawalHistoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RakutenSecurityScraperServer).ListWithdrawalHistories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RakutenSecurityScraper_ListWithdrawalHistories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RakutenSecurityScraperServer).ListWithdrawalHistories(ctx, req.(*ListWithdrawalHistoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RakutenSecurityScraper_ListDividendHistories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDividendHistoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RakutenSecurityScraperServer).ListDividendHistories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RakutenSecurityScraper_ListDividendHistories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RakutenSecurityScraperServer).ListDividendHistories(ctx, req.(*ListDividendHistoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RakutenSecurityScraper_TotalAssets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TotalAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RakutenSecurityScraperServer).TotalAssets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RakutenSecurityScraper_TotalAssets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RakutenSecurityScraperServer).TotalAssets(ctx, req.(*TotalAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RakutenSecurityScraper_ServiceDesc is the grpc.ServiceDesc for RakutenSecurityScraper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RakutenSecurityScraper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RakutenSecurityScraper",
	HandlerType: (*RakutenSecurityScraperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListWithdrawalHistories",
			Handler:    _RakutenSecurityScraper_ListWithdrawalHistories_Handler,
		},
		{
			MethodName: "ListDividendHistories",
			Handler:    _RakutenSecurityScraper_ListDividendHistories_Handler,
		},
		{
			MethodName: "TotalAssets",
			Handler:    _RakutenSecurityScraper_TotalAssets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rakuten-security-scraper.proto",
}
