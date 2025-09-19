package flow

import (
	"context"
	"fmt"
	"go-tpl/server/infra/global"
	"testing"
)

func TestBiddingEngine_BiddingMatch(t *testing.T) {
	store := &RedisStorage{
		Redis: global.Redis,
	}

	e := &BiddingEngine{
		refind: &RefindFilter{},
		retrieval: &CompositeRetrieval{
			retrievals: []Retrieval{
				&CityRetrieval{store: store},
				&CinemaRetrieval{store: store},
				&MovieHallRetrieval{store: store},
				&MovieNameRetrieval{store: store},
			},
			store: store,
		},
	}

	params := BiddingParamsDTO{
		City:      "北京",
		MovieName: "哪吒2",
	}
	match, err := e.MatchBidding(context.Background(), params, 3)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	fmt.Println(match)
}
