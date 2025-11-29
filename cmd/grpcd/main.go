package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"signalstack/internal/data/storage"
	"signalstack/internal/engine"
	"signalstack/internal/grpcserver"
	"signalstack/internal/service"
	"signalstack/internal/strategy"
	signalstackv1 "signalstack/proto"
)

func main() {
    // 1) Wire dependencies
    store := storage.NewCSVStore("./data/raw/yahoo")

    reg := strategy.NewRegistry()
    strategy.RegisterBuiltins(reg)

    eng := engine.New(engine.EngineConfig{
        InitialCapital:     10000,
        SlippageBps:        5,
        CommissionPerTrade: 0,
        AllInOnEntry:       true,
    })

    svc := service.NewBacktestService(store, reg, eng)
    grpcSrv := grpcserver.NewBacktestGRPCServer(svc)

    // 2) Start gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    signalstackv1.RegisterBacktestServiceServer(s, grpcSrv)

    // optional: enable server reflection for grpcurl, etc.
    reflection.Register(s)

    log.Println("Backtest gRPC server listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}