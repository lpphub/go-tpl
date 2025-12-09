package permission

import (
	"context"
	"errors"
	"go-tpl/logic/role"
	"go-tpl/logic/shared"
	"go-tpl/web/types"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// List 获取权限列表
func (s *Service) List(ctx context.Context, req types.PermissionQueryReq) (*shared.PageData[Permission], error) {
	var (
		total int64
		list  []Permission
	)

	_db := s.db.WithContext(ctx).Model(&Permission{})
	if req.Code != "" {
		_db = _db.Where("code like ?", "%"+req.Code+"%")
	}
	if req.Name != "" {
		_db = _db.Where("name like ?", "%"+req.Name+"%")
	}
	if req.Module != "" {
		_db = _db.Where("module = ?", req.Module)
	}
	if req.Status != nil {
		_db = _db.Where("status = ?", *req.Status)
	}

	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		if err := _db.Scopes(shared.Paginate(req.Pagination)).Find(&list).Error; err != nil {
			return nil, err
		}
	}
	return shared.WithPageData[Permission](total, list), nil
}

// Get 获取单个权限
func (s *Service) Get(ctx context.Context, id uint) (*Permission, error) {
	var permission Permission
	if err := s.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrRecordNotFound
		}
		return nil, err
	}
	return &permission, nil
}

// Create 创建权限
func (s *Service) Create(ctx context.Context, req types.CreatePermissionReq) (*Permission, error) {
	// 检查权限代码是否已存在
	var count int64
	if err := s.db.WithContext(ctx).Model(&Permission{}).Where("code = ?", req.Code).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, shared.ErrPermissionExists
	}

	permission := Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Module:      req.Module,
		Status:      shared.StatusActive,
	}

	if err := s.db.WithContext(ctx).Create(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

// Update 更新权限
func (s *Service) Update(ctx context.Context, id uint, req types.UpdatePermissionReq) error {
	permission, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.Code != "" && req.Code != permission.Code {
		// 检查权限代码是否已存在
		var count int64
		if err = s.db.WithContext(ctx).Model(&Permission{}).Where("code = ? AND id != ?", req.Code, id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return shared.ErrPermissionExists
		}
		updates["code"] = req.Code
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Module != nil {
		updates["module"] = *req.Module
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return shared.ErrInvalidParam
	}

	return s.db.WithContext(ctx).Model(&Permission{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除权限
func (s *Service) Delete(ctx context.Context, id uint) error {
	// 检查权限是否存在
	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// 开启事务，同时删除角色权限关联
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联
		if err = tx.Where("permission_id = ?", id).Delete(&role.RolePermission{}).Error; err != nil {
			return err
		}

		// 删除权限
		return tx.Delete(&Permission{}, id).Error
	})
}

// UpdateStatus 更新权限状态
func (s *Service) UpdateStatus(ctx context.Context, id uint, status int8) error {
	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&Permission{}).Where("id = ?", id).Update("status", status).Error
}

// GetModules 获取所有模块列表
func (s *Service) GetModules(ctx context.Context) ([]string, error) {
	var modules []string
	err := s.db.WithContext(ctx).Model(&Permission{}).Distinct("module").Pluck("module", &modules).Error
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// GetPermissionRoles 获取权限角色列表
func (s *Service) GetPermissionRoles(ctx context.Context, permissionId uint) ([]uint, error) {
	var rolePermissions []role.RolePermission
	err := s.db.WithContext(ctx).Where("permission_id = ?", permissionId).Find(&rolePermissions).Error
	if err != nil {
		return nil, err
	}

	roleIds := make([]uint, len(rolePermissions))
	for i, rp := range rolePermissions {
		roleIds[i] = rp.RoleID
	}
	return roleIds, nil
}
