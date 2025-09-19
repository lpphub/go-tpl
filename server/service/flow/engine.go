package flow

import (
	"context"
)

type Engine[T, R any] interface {
	Recall(ctx context.Context, t T) (R, error)
	PreRanking(ctx context.Context, t T) (R, error)
	Ranking(ctx context.Context, t T, topN int) (R, error)
}

// BiddingEngine 竞价引擎
/**
 * 竞价引擎
 * 1. retrieval：基于redis按部分维度进行数据召回，缩小范围
 * 2. refind：基于召回数据详情进行更精细化的过滤
 */
type BiddingEngine struct {
	retrieval *CompositeRetrieval
	refind    *RefindFilter
	startKey  string
}

func (e *BiddingEngine) MatchBidding(ctx context.Context, params BiddingParamsDTO, topN int) ([]*BiddingResultDTO, error) {
	// 1.基于Redis的城市、影院、影厅、电影匹配
	tempKey, err := e.retrieval.Recall(ctx, e.startKey, params)
	if err != nil {
		return nil, err
	}

	// 2.获取召回的规则ID
	ruleIds, err := e.retrieval.GetRecallIds(ctx, tempKey)
	if err != nil || len(ruleIds) == 0 {
		return nil, err
	}

	// 3.精细化过滤
	refindData, err := e.refind.Filter(ctx, ruleIds, params)
	if err != nil {
		return nil, err
	}

	// 4.竞价ranking
	return e.ranking(ctx, refindData, topN)
}

func (e *BiddingEngine) ranking(_ context.Context, data []*BiddingResultDTO, topN int) ([]*BiddingResultDTO, error) {
	// 排序

	// TopN结果
	if len(data) > topN {
		return data[:topN], nil
	}
	return data, nil
}
