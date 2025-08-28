package flow

import (
	"context"
)

type Retrieval interface {
	Recall(ctx context.Context, key string, params BiddingParamsDTO) (string, error)
}

type CityRetrieval struct {
	store Storage
}

func (c CityRetrieval) Recall(ctx context.Context, key string, params BiddingParamsDTO) (string, error) {
	//TODO implement me
	panic("implement me")
}

type CinemaRetrieval struct {
	store Storage
}

func (c CinemaRetrieval) Recall(ctx context.Context, key string, params BiddingParamsDTO) (string, error) {
	//TODO implement me
	panic("implement me")
}

type MovieHallRetrieval struct {
	store Storage
}

func (m MovieHallRetrieval) Recall(ctx context.Context, key string, params BiddingParamsDTO) (string, error) {
	//TODO implement me
	panic("implement me")
}

type MovieNameRetrieval struct {
	store Storage
}

func (m MovieNameRetrieval) Recall(ctx context.Context, key string, params BiddingParamsDTO) (string, error) {
	//TODO implement me
	panic("implement me")
}

type CompositeRetrieval struct {
	store      Storage
	retrievals []Retrieval
}

func (c *CompositeRetrieval) Recall(ctx context.Context, key string, params BiddingParamsDTO) (string, error) {
	tempKey := key
	for _, retrieval := range c.retrievals {
		var err error
		tempKey, err = retrieval.Recall(ctx, tempKey, params)
		if err != nil {
			return "", err
		}
		count, err := c.store.GetCount(ctx, tempKey)
		if err != nil || count == 0 {
			return tempKey, nil
		}
	}
	return tempKey, nil
}

func (c *CompositeRetrieval) GetRecallIds(ctx context.Context, setKey string) ([]string, error) {
	return c.store.GetMembers(ctx, setKey)
}
