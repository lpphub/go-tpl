package entity

import "github.com/shopspring/decimal"

type BiddingOtherRule struct {
	Id                 int64           `gorm:"column:id;primaryKey"`         // 主键ID（数据库为 int unsigned，对应Go的uint64）
	UserId             int64           `gorm:"column:user_id"`               // 关联用户ID（票商ID）
	Name               string          `gorm:"column:name"`                  // 规则名称
	CityName           string          `gorm:"column:city_name"`             // 包含城市（多城市用字符串存储，如逗号分隔）
	CityNotIn          string          `gorm:"column:city_not_in"`           // 不包含城市（多城市用字符串存储）
	CinemaBrandId      string          `gorm:"column:cinema_brand_id"`       // 影院品牌ID（多品牌用字符串存储）
	CinemaCategoryId   string          `gorm:"column:cinema_category_id"`    // 自定义影院类目ID（包含）
	NoCinemaCategoryId string          `gorm:"column:no_cinema_category_id"` // 自定义影院类目ID（不包含）
	RuleType           int8            `gorm:"column:rule_type"`             // 规则模式：1.连锁模式，2.杂牌模式，3.电影模式
	UpdateDate         string          `gorm:"column:update_date"`           // 规则更新时间（自动更新datetime）
	Status             int8            `gorm:"column:status"`                // 规则状态：1.开启，2.禁用
	CinemaIds          string          `gorm:"column:cinema_ids"`            // 不包含影院ID（多ID用文本存储，如逗号分隔）
	CinemaIdsIn        string          `gorm:"column:cinema_ids_in"`         // 包含影院ID（多ID用文本存储）
	HallIds            string          `gorm:"column:hall_ids"`              // 影厅分类ID（多ID用字符串存储）
	MovieName          string          `gorm:"column:movie_name"`            // 包含电影名称（多电影用字符串存储）
	Discount           decimal.Decimal `gorm:"column:discount"`              // 折扣比例（如9.5表示95折）
	CreateTime         int             `gorm:"column:create_time"`           // 规则创建时间（时间戳）
	LoverSeat          int8            `gorm:"column:lover_seat"`            // 情侣座配置：0不包含，1包含
	AdminId            int64           `gorm:"column:admin_id"`              // 规则添加管理员ID：0为票商自己添加
	UpdateTime         string          `gorm:"column:update_time"`           // 规则更新时间（与update_date功能一致，冗余字段）
	IfBottomCover      int8            `gorm:"column:ifBottomCover"`         // 是否保底报价：0否，1是
	MovieType          string          `gorm:"column:movie_type"`            // 电影类型（如"2D"、"3D"等）
}

func (BiddingOtherRule) TableName() string {
	return "tb_hh_auto_bidding_rule"
}
