package checker

import (
	"go-tpl/internal/domain/repo"
	"go-tpl/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
)

type LocalCacheLoader struct {
	brandRepo    *repo.CinemaBrandRepo
	categoryRepo *repo.CustomCinemaCategoryRepo
}

var (
	cacheBrand    = cache.New(10*time.Minute, 15*time.Minute)
	cacheCategory = cache.New(5*time.Minute, 10*time.Minute)
)

func NewLocalCacheLoader(brandRepo *repo.CinemaBrandRepo, categoryRepo *repo.CustomCinemaCategoryRepo) *LocalCacheLoader {
	return &LocalCacheLoader{
		brandRepo:    brandRepo,
		categoryRepo: categoryRepo,
	}
}

func (l *LocalCacheLoader) PreloadCache(ctx *gin.Context) error {
	if _, exist := cacheBrand.Get("brand:exist"); !exist {
		// 1.缓存品牌对应的影院
		allBrand := l.brandRepo.ListAll(ctx)
		for _, brand := range allBrand {
			cacheBrand.Set(cast.ToString(brand.Id), brand.CinemaIds, cache.DefaultExpiration)
		}
		cacheBrand.Set("brand:exist", "1", cache.DefaultExpiration)
	}

	if _, exist := cacheCategory.Get("category:exist"); !exist {
		// 2. 查询所有影院与分类关系
		allCategory := l.categoryRepo.ListAll(ctx)
		for _, category := range allCategory {
			cacheCategory.Set(cast.ToString(category.Id), category.CinemaId, cache.DefaultExpiration)
		}
		cacheCategory.Set("category:exist", "1", cache.DefaultExpiration)
	}
	return nil
}

func getCinemaIdsByBrandId(brandId int64) []int64 {
	if v, ok := cacheBrand.Get(cast.ToString(brandId)); ok {
		cinemaIds := v.(string)
		return util.SplitToInt64Slice(cinemaIds, ",")
	}
	return []int64{}
}

func getCinemaIdsByBrandIds(brandIds []int64) []int64 {
	var cinemaIds []int64
	for _, id := range brandIds {
		cinemaId := getCinemaIdsByBrandId(id)
		cinemaIds = append(cinemaIds, cinemaId...)
	}
	return cinemaIds
}

func getCinemaIdsByCategoryId(categoryId int64) []int64 {
	if v, ok := cacheCategory.Get(cast.ToString(categoryId)); ok {
		cinemaIds := v.(string)
		return util.SplitToInt64Slice(cinemaIds, ",")
	}
	return []int64{}
}

func getCinemaIdsByCategoryIds(categoryIds []int64) []int64 {
	var cinemaIds []int64
	for _, id := range categoryIds {
		cinemaId := getCinemaIdsByCategoryId(id)
		cinemaIds = append(cinemaIds, cinemaId...)
	}
	return cinemaIds
}
