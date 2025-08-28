package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BiddingBillRepo struct {
	db *gorm.DB
}

func NewBiddingBillRepo(db *gorm.DB) *BiddingBillRepo {
	return &BiddingBillRepo{db: db}
}

func (b *BiddingBillRepo) InsertWithTx(tx *gorm.DB, bill *entity.BiddingBill) error {
	return tx.Create(bill).Error
}

func (b *BiddingBillRepo) ListByUidAndNotOrderId(ctx *gin.Context, uid, orderId int64) ([]entity.BiddingBill, error) {
	var list []entity.BiddingBill
	err := b.db.WithContext(ctx).Model(&entity.BiddingBill{}).Where("uid = ? and order_id != ?", uid, orderId).
		Where("status in ? and delete_time = 0", []int8{0, 1, 3}).
		Select([]string{"price", "uid", "order_id", "is_auto"}).
		Scan(&list).Error
	return list, err
}
