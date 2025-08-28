package repo

import (
	"go-tpl/internal/domain/entity"
	"go-tpl/pkg/pagination"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetOne(ctx *gin.Context, uid int64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", uid).First(&user).Error
	return &user, err
}

func (r *UserRepo) PageList(ctx *gin.Context, paginate pagination.Pagination) (*pagination.PageData[entity.User], error) {
	var (
		total int64
		list  []entity.User
	)
	_db := r.db.WithContext(ctx).Model(&entity.User{})

	//if param.Uid > 0 {
	//	_db.Where("id = ?", param.Uid)
	//}

	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		_db.Order("id desc").Scopes(Paginate(paginate.Pn, paginate.Ps)).Find(&list)
		return pagination.WrapPageData(total, list), nil
	}
	return pagination.WrapPageData(total, Empty[entity.User]()), nil
}

func (r *UserRepo) UpdateWithTx(tx *gorm.DB, uid int64, updates map[string]any) error {
	return tx.Model(&entity.User{}).Where("id = ?", uid).Updates(updates).Error
}
