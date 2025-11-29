package grpcserver

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"signalstack/internal/service"
	"signalstack/pkg/api"
	signalstackv1 "signalstack/proto"
)

type BacktestGRPCServer struct {
	signalstackv1.UnimplementedBacktestServiceServer
	svc *service.BacktestService
}

func NewBacktestGRPCServer(svc *service.BacktestService) *BacktestGRPCServer {
	return &BacktestGRPCServer{
		svc: svc,
	}
}

func (s *BacktestGRPCServer) RunBacktest(
	ctx context.Context,
	req *signalstackv1.RunBacktestRequest,
) (*signalstackv1.RunBacktestResponse, error) {
	// convert proto to internal api request
	start := fromTimestamp(req.GetStart())
	end := fromTimestamp(req.GetEnd())
	apiReq := api.BacktestRequest{
		JobID:     req.GetJobId(),
		Symbol:    req.GetSymbol(),
		Timeframe: req.GetTimeframe(),
		StartDate: start,
		EndDate:   end,
		Strategy:  req.GetStrategyId(),
		Params:    req.GetParams(),
	}

	// call the backtest service
	apiResp, err := s.svc.Run(apiReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to run backtest: %v", err)
	}

	// convert trades
	pbTrades := make([]*signalstackv1.Trade, 0, len(apiResp.Trades))
	for _, trade := range apiResp.Trades {
		r := 0.0
		if trade.RMultiple != nil {
			r = *trade.RMultiple
		}
		pbTrades = append(pbTrades, &signalstackv1.Trade{
			Side:       trade.Side,
			Position:   trade.Position,
			EntryTime:  timestamppb.New(trade.EntryTime),
			EntryPrice: trade.Entry,
			ExitTime:   timestamppb.New(trade.ExitTime),
			ExitPrice:  trade.Exit,
			Size:       trade.Size,
			Pnl:        trade.PnL,
			RMultiple:  r,
		})
	}

	// convert api response to proto response
	protoResp := &signalstackv1.RunBacktestResponse{
		JobId:      apiResp.JobID,
		Symbol:     apiResp.Symbol,
		StrategyId: apiResp.Strategy,
		Metrics: &signalstackv1.BacktestMetrics{
			Cagr:         apiResp.Metrics.CAGR,
			MaxDrawdown:  apiResp.Metrics.MaxDrawdown,
			WinRate:      apiResp.Metrics.WinRate,
			ProfitFactor: apiResp.Metrics.ProfitFactor,
			FinalEquity:  apiResp.Metrics.FinalEquity,
		},
		Trades: pbTrades,
	}
	return protoResp, nil
}

func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
