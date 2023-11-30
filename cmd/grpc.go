package main

import (
	"context"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/irisco88/authutil"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

func grpcServer(logger *zap.Logger, logReqs bool) *grpc.Server {
	streamServerOptions := []grpc.StreamServerInterceptor{
		grpcTags.StreamServerInterceptor(grpcTags.WithFieldExtractor(grpcTags.CodeGenRequestFieldExtractor)),
		grpcZap.StreamServerInterceptor(logger),
		grpcZap.PayloadStreamServerInterceptor(logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return logReqs
		}),
		grpcRecovery.StreamServerInterceptor(grpcRecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("stack trace from panic " + string(debug.Stack()))
			return status.Errorf(codes.Internal, "%v", p)
		})),
		authutil.StreamServerInterceptor(),
	}

	unaryServerOptions := []grpc.UnaryServerInterceptor{
		grpcTags.UnaryServerInterceptor(grpcTags.WithFieldExtractor(grpcTags.CodeGenRequestFieldExtractor)),
		grpcZap.UnaryServerInterceptor(logger),
		grpcZap.PayloadUnaryServerInterceptor(logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return logReqs
		}),
		grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("stack trace from panic " + string(debug.Stack()))
			return status.Errorf(codes.Internal, "%v", p)
		})),
		authutil.UnaryServerInterceptor(),
	}
	return grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(streamServerOptions...)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(unaryServerOptions...)),
		grpc.MaxRecvMsgSize(1024*1024*50),
	)
}
