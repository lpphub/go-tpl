package entity

import (
	"go-tpl/pkg/ext"

	"github.com/shopspring/decimal"
)

type BiddingCustomRule struct {
	Id                 int64           `gorm:"column:id;primaryKey"`
	Uid                int64           `gorm:"column:uid"`
	Name               string          `gorm:"column:name"`
	Province           string          `gorm:"column:province"`
	City               string          `gorm:"column:city"`
	Platform           string          `gorm:"column:platform"`
	HallName           string          `gorm:"column:hallname"`
	Movie              string          `gorm:"column:movie"`
	NoMovie            string          `gorm:"column:no_movie"`
	NoHallName         string          `gorm:"column:no_hallname"`
	NoPlatform         string          `gorm:"column:no_platform"`
	CinemaCategoryId   string          `gorm:"column:cinema_category_id"`
	NoCinemaCategoryId string          `gorm:"column:no_cinema_category_id"`
	NoCity             string          `gorm:"column:no_city"`
	NoProvince         string          `gorm:"column:no_province"`
	PriceMin           decimal.Decimal `gorm:"column:price_min"`
	PriceMax           decimal.Decimal `gorm:"column:price_max"`
	Num                string          `gorm:"column:num"`
	LoverSeat          int8            `gorm:"column:lover_seat"`
	IsSeat             int8            `gorm:"column:is_seat"`
	Type               int8            `gorm:"column:type"`
	Value              decimal.Decimal `gorm:"column:value"`
	Status             int8            `gorm:"column:status"`
	CreateAt           int             `gorm:"column:create_time"`
	UpdateAt           *ext.Timestamp  `gorm:"column:update_time"`
	DeleteAt           int             `gorm:"column:delete_time"`
	CinemaId           string          `gorm:"column:cinemaId"`
	NoCinemaId         string          `gorm:"column:noCinemaId"`
	MovieName          string          `gorm:"column:movieName"`
	NoMovieName        string          `gorm:"column:noMovieName"`
	CityName           string          `gorm:"column:cityName"`
	NoCityName         string          `gorm:"column:noCityName"`
	BrandId            string          `gorm:"column:brandId"`
	HallIds            string          `gorm:"column:hallIds"`
	NoHallIds          string          `gorm:"column:noHallIds"`
	AdminId            int64           `gorm:"column:admin_id"`
	UpdateTimeSecond   string          `gorm:"column:update_time"`
	IfBottomCover      int8            `gorm:"column:ifBottomCover"`
	HourSpan           string          `gorm:"column:hourSpan"`
	WeekNum            string          `gorm:"column:weekNum"`
	MovieType          string          `gorm:"column:movie_type"`
}

func (*BiddingCustomRule) TableName() string {
	return "tb_bidding_auto"
}
