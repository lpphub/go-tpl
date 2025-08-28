package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomCinemaCategoryRepo struct {
	db *gorm.DB
}

func NewCustomCinemaCategoryRepo(db *gorm.DB) *CustomCinemaCategoryRepo {
	return &CustomCinemaCategoryRepo{db: db}
}

func (r *CustomCinemaCategoryRepo) ListByIds(ctx *gin.Context, ids []int64) ([]*entity.UserCustomCinemaCategory, error) {
	var list []*entity.UserCustomCinemaCategory
	err := r.db.WithContext(ctx).Model(&entity.UserCustomCinemaCategory{}).Where("id in ?", ids).Scan(&list).Error
	return list, err
}

func (r *CustomCinemaCategoryRepo) ListAll(ctx *gin.Context) []*entity.UserCustomCinemaCategory {
	var list []*entity.UserCustomCinemaCategory
	r.db.WithContext(ctx).Model(&entity.UserCustomCinemaCategory{}).Where("status = ?", 1).Scan(&list)
	return list
}
