package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OupAppRepo struct {
	db *gorm.DB
}

func NewOupAppRepo(db *gorm.DB) *OupAppRepo {
	return &OupAppRepo{db: db}
}

func (o *OupAppRepo) Exist(ctx *gin.Context, uid int64) bool {
	var count int64
	o.db.WithContext(ctx).Model(&entity.OupAppSecret{}).Where("userId = ?", uid).Count(&count)
	return count > 0
}
