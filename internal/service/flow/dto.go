package flow

import "github.com/shopspring/decimal"

type BiddingParamsDTO struct {
	City               string // 城市
	Cinema             string // 影院
	MovieHall          string // 影厅
	MovieName          string // 电影名称
	MovieType          int8   // 电影类型 2D 3D
	MovieShowTimeStart int    // 电影放映开始时间
	MovieShowTimeEnd   int    // 电影放映结束时间
	MovieHasWeek       string // 电影包含的星期
	LoverSeat          int8   // 是否情侣座
	ChangeSeat         int8   // 是否可换座
	SeatNum            int    // 座位数（票数）
}

type BiddingResultDTO struct {
	Status     int8
	RuleId     uint64
	UID        uint64
	Price      decimal.Decimal
	FailReason string
}
