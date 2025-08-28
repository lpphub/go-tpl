package repo

import (
	"go-tpl/internal/domain/entity"

	"gorm.io/gorm"
)

type BiddingBillBottomRepo struct {
	db *gorm.DB
}

func NewBiddingBillBottomRepo(db *gorm.DB) *BiddingBillBottomRepo {
	return &BiddingBillBottomRepo{db: db}
}

func (r *BiddingBillBottomRepo) InsertWithTx(tx *gorm.DB, bill *entity.BiddingBillBottomCover) error {
	return tx.Create(bill).Error
}
