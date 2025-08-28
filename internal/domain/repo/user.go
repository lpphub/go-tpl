package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetOne(ctx *gin.Context, userId int64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&user).Error
	return &user, err
}

func (r *UserRepo) ListEffectiveByIds(ctx *gin.Context, userIds []int64) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).Where("user_id in ?", userIds).
		Where("bidding_status = 1 and is_ok = 1 and bid_points >= 150 and status != 1").Find(&users).Error
	return users, err
}

func (r *UserRepo) UpdateWithTx(tx *gorm.DB, userId int64, updates map[string]any) error {
	return tx.Model(&entity.User{}).Where("user_id = ?", userId).Updates(updates).Error
}
