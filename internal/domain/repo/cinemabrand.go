package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CinemaBrandRepo struct {
	db *gorm.DB
}

func NewCinemaBrandRepo(db *gorm.DB) *CinemaBrandRepo {
	return &CinemaBrandRepo{db: db}
}

func (r *CinemaBrandRepo) GetOne(ctx *gin.Context, id int64) (*entity.CinemaBrand, error) {
	var brand entity.CinemaBrand
	err := r.db.WithContext(ctx).First(&brand, id).Error
	return &brand, err
}

func (r *CinemaBrandRepo) ListByBrandIds(ctx *gin.Context, ids []int64) ([]*entity.CinemaBrand, error) {
	var brands []*entity.CinemaBrand
	err := r.db.WithContext(ctx).Where("id in ?", ids).Find(&brands).Error
	return brands, err
}

func (r *CinemaBrandRepo) ListAll(ctx *gin.Context) []*entity.CinemaBrand {
	var brands []*entity.CinemaBrand
	r.db.WithContext(ctx).Model(&entity.CinemaBrand{}).Where("status = ?", 1).Find(&brands)
	return brands
}
