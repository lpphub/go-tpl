package role

import (
	"context"
	"errors"
	"go-tpl/logic/shared"
	"go-tpl/logic/user"
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

// List 获取角色列表
func (s *Service) List(ctx context.Context, req types.RoleQueryReq) (*shared.PageData[Role], error) {
	var (
		total int64
		list  []Role
	)

	_db := s.db.WithContext(ctx).Model(&Role{})
	if req.Name != "" {
		_db = _db.Where("name like ?", "%"+req.Name+"%")
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
	return shared.Wrapper[Role](total, list), nil
}

// Get 获取单个角色
func (s *Service) Get(ctx context.Context, id uint) (*Role, error) {
	var role Role
	if err := s.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrRecordNotFound
		}
		return nil, err
	}
	return &role, nil
}

// Create 创建角色
func (s *Service) Create(ctx context.Context, req types.CreateRoleReq) (*Role, error) {
	// 检查角色名是否已存在
	var count int64
	if err := s.db.WithContext(ctx).Model(&Role{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, shared.ErrRoleExists
	}

	role := Role{
		Name:        req.Name,
		Description: req.Description,
		Status:      shared.StatusActive,
	}

	if err := s.db.WithContext(ctx).Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// Update 更新角色
func (s *Service) Update(ctx context.Context, id uint, req types.UpdateRoleReq) error {
	role, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.Name != "" && req.Name != role.Name {
		// 检查角色名是否已存在
		var count int64
		if err := s.db.WithContext(ctx).Model(&Role{}).Where("name = ? AND id != ?", req.Name, id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return shared.ErrRoleExists
		}
		updates["name"] = req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return shared.ErrInvalidParam
	}

	return s.db.WithContext(ctx).Model(&Role{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除角色
func (s *Service) Delete(ctx context.Context, id uint) error {
	// 检查角色是否存在
	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// 开启事务，同时删除角色权限关联
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联
		if err := tx.Where("role_id = ?", id).Delete(&RolePermission{}).Error; err != nil {
			return err
		}

		// 删除用户角色关联
		if err := tx.Where("role_id = ?", id).Delete(&user.UserRole{}).Error; err != nil {
			return err
		}

		// 删除角色
		return tx.Delete(&Role{}, id).Error
	})
}

// UpdateStatus 更新角色状态
func (s *Service) UpdateStatus(ctx context.Context, id uint, status int8) error {
	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&Role{}).Where("id = ?", id).Update("status", status).Error
}

// GetRolePermissions 获取角色权限列表
func (s *Service) GetRolePermissions(ctx context.Context, roleId uint) ([]uint, error) {
	var rolePermissions []RolePermission
	err := s.db.WithContext(ctx).Where("role_id = ?", roleId).Find(&rolePermissions).Error
	if err != nil {
		return nil, err
	}

	permissionIds := make([]uint, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissionIds[i] = rp.PermissionID
	}
	return permissionIds, nil
}

// AssignPermissions 为角色分配权限
func (s *Service) AssignPermissions(ctx context.Context, roleId uint, permissionIds []uint) error {
	// 检查角色是否存在
	_, err := s.Get(ctx, roleId)
	if err != nil {
		return err
	}

	// 开启事务
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色现有权限
		if err := tx.Where("role_id = ?", roleId).Delete(&RolePermission{}).Error; err != nil {
			return err
		}

		// 分配新权限
		if len(permissionIds) > 0 {
			rolePermissions := make([]RolePermission, len(permissionIds))
			for i, permissionId := range permissionIds {
				rolePermissions[i] = RolePermission{
					RoleID:       roleId,
					PermissionID: permissionId,
				}
			}
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetRoleUsers 获取角色用户列表
func (s *Service) GetRoleUsers(ctx context.Context, roleId uint) ([]uint, error) {
	var userRoles []user.UserRole
	err := s.db.WithContext(ctx).Where("role_id = ?", roleId).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	userIds := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		userIds[i] = ur.UserID
	}
	return userIds, nil
}
