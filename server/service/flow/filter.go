package flow

import (
	"context"
)

type Filter interface {
	Filter(ctx context.Context, ruleIds []string, params BiddingParamsDTO) ([]*BiddingResultDTO, error)
}

type RefindFilter struct {
}

func (r *RefindFilter) Filter(ctx context.Context, ruleIds []string, params BiddingParamsDTO) ([]*BiddingResultDTO, error) {
	return nil, nil
}
