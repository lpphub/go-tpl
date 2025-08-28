package dto

import (
	"go-tpl/internal/domain/entity"
)

type OrderDTO struct {
	entity.Order
	EndTime         int64
	HallClassifyId  int64  // 影厅ID（单个）
	HallClassifyIds string // 影厅ID（多个匹配模式，字符串存储）
	MovieType       string // 电影制式（如"2D"、"3D"、"IMAX"）
	IsBottomCover   bool
}

func (o *OrderDTO) GetIsBottomCover() bool {
	if o.LimitTime > 0 {
		timeDiff := o.LimitTime - o.CreateTime

		if timeDiff < 0 {
			timeDiff = 0
		}

		if o.IsBottomCover && timeDiff <= 360 && timeDiff > 0 {
			return true
		}
	}
	return false
}
