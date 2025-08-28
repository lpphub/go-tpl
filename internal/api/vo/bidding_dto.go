package vo

type BiddingAutoDTO struct {
	OrderId            int64 `json:"orderId" binding:"required"`
	IsBidding          bool  `json:"ifDoBidding"`
	IsLimitPrice       bool  `json:"ifLimitPrice"`
	IsSaveRedis        bool  `json:"saveRedis"`
	OrderLimitTime     int   `json:"orderLimitTime"`
	IsOrderBottomCover bool  `json:"ifOrderIsBottomCover"`
}
