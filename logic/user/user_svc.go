package user

import (
	"context"
	"errors"
	"go-tpl/logic/shared"
	"go-tpl/web/types"

	"golang.org/x/crypto/bcrypt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewService(db *gorm.DB, redis *redis.Client) *Service {
	return &Service{
		db:    db,
		redis: redis,
	}
}

// List 获取用户列表
func (s *Service) List(ctx context.Context, req types.UserQueryReq) (*shared.PageData[User], error) {
	var (
		total int64
		list  []User
	)

	_db := s.db.WithContext(ctx).Model(&User{})
	if req.Username != "" {
		_db = _db.Where("username like ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		_db = _db.Where("email like ?", "%"+req.Email+"%")
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
	return shared.Wrapper[User](total, list), nil
}

// Get 获取单个用户
func (s *Service) Get(ctx context.Context, id uint) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (s *Service) Create(ctx context.Context, req types.CreateUserReq) (*User, error) {
	// 检查用户名是否已存在
	var count int64
	if err := s.db.WithContext(ctx).Model(&User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, shared.ErrUserExists
	}

	// 检查邮箱是否已存在
	if err := s.db.WithContext(ctx).Model(&User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, shared.ErrEmailExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, shared.ErrInvalidPassword
	}

	user := User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Status:   shared.StatusActive,
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Update 更新用户
func (s *Service) Update(ctx context.Context, id uint, req types.UpdateUserReq) error {
	user, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.Username != "" && req.Username != user.Username {
		// 检查用户名是否已存在
		var count int64
		if err := s.db.WithContext(ctx).Model(&User{}).Where("username = ? AND id != ?", req.Username, id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return shared.ErrUserExists
		}
		updates["username"] = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		// 检查邮箱是否已存在
		var count int64
		if err := s.db.WithContext(ctx).Model(&User{}).Where("email = ? AND id != ?", req.Email, id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return shared.ErrEmailExists
		}
		updates["email"] = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return shared.ErrInvalidPassword
		}
		updates["password"] = string(hashedPassword)
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return shared.ErrInvalidParam
	}

	return s.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除用户
func (s *Service) Delete(ctx context.Context, id uint) error {
	// 检查用户是否存在
	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Delete(&User{}, id).Error
}

// UpdateStatus 更新用户状态
func (s *Service) UpdateStatus(ctx context.Context, id uint, status int8) error {
	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("status", status).Error
}

// GetUserRoles 获取用户角色列表
func (s *Service) GetUserRoles(ctx context.Context, userId uint) ([]uint, error) {
	var userRoles []UserRole
	err := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	roleIds := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		roleIds[i] = ur.RoleID
	}
	return roleIds, nil
}

// AssignRoles 为用户分配角色
func (s *Service) AssignRoles(ctx context.Context, userId uint, roleIds []uint) error {
	// 检查用户是否存在
	_, err := s.Get(ctx, userId)
	if err != nil {
		return err
	}

	// 开启事务
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除用户现有角色
		if err := tx.Where("user_id = ?", userId).Delete(&UserRole{}).Error; err != nil {
			return err
		}

		// 分配新角色
		if len(roleIds) > 0 {
			userRoles := make([]UserRole, len(roleIds))
			for i, roleId := range roleIds {
				userRoles[i] = UserRole{
					UserID: userId,
					RoleID: roleId,
				}
			}
			if err := tx.Create(&userRoles).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ValidateLogin 验证用户登录
func (s *Service) ValidateLogin(ctx context.Context, username, password string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewError(2104, "用户名或密码错误")
		}
		return nil, err
	}

	// 检查用户状态
	if user.Status != shared.StatusActive {
		return nil, shared.NewError(2105, "用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, shared.NewError(2104, "用户名或密码错误")
	}

	return &user, nil
}
