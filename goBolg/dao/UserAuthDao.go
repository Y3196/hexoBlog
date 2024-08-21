package dao

import (
	"context"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
)

// UserAuthDao 定义了操作用户认证数据的接口
type UserAuthDao interface {
	ListUsers(ctx context.Context, current, size int, condition vo.ConditionVO) ([]dto.UserBackDTO, error)
	CountUser(ctx context.Context, condition vo.ConditionVO) (int64, error)
	UpdateUserAuth(ctx context.Context, userAuth *model.UserAuth) error
	InsertUserAuth(ctx context.Context, userAuth *model.UserAuth) error
	UpdatePassword(ctx context.Context, username, hashedPassword string) error
	FindUserById(ctx context.Context, id int) (*model.UserAuth, error)
	FindUserByUsername(ctx context.Context, username string) (*model.UserAuth, error)
	FindIpSourceById(ctx context.Context, id uint) (string, error)
	SelectUserByUsername(ctx context.Context, username string) (*model.UserAuth, error)
	UpdateById(ctx context.Context, userAuth *model.UserAuth) error
}

// uniqueViewDao 实现 UserAuthDao 接口
type userAuthDao struct {
	db *gorm.DB
}

// NewUserAuthDao 创建一个新的 UserAuthDao 实例
func NewUserAuthDao(db *gorm.DB) UserAuthDao {
	return &userAuthDao{db: db}
}

// ListUsers 查询后台用户列表
func (dao *userAuthDao) ListUsers(ctx context.Context, current, size int, condition vo.ConditionVO) ([]dto.UserBackDTO, error) {
	var users []dto.UserBackDTO
	query := dao.db.WithContext(ctx).
		Table("tb_user_info").
		Select("ui.id, ui.avatar, ui.nickname, ui.is_disable, ua.user_info_id, ua.avatar, ua.nickname, ua.login_type, r.id as role_id, r.role_name, ua.ip_address, ua.ip_source, ua.create_time, ua.last_login_time").
		Joins("LEFT JOIN tb_user_auth ua ON ua.user_info_id = ui.id").
		Joins("LEFT JOIN tb_user_role ur ON ui.id = ur.user_id").
		Joins("LEFT JOIN tb_role r ON ur.role_id = r.id").
		Where("1 = 1")

	if condition.LoginType != nil {
		query = query.Where("ua.login_type = ?", *condition.LoginType)
	}
	if condition.Keywords != nil {
		query = query.Where("ui.nickname LIKE ?", "%"+*condition.Keywords+"%")
	}

	err := query.Offset(current).Limit(size).Order("ui.id").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CountUser 查询后台用户数量
func (dao *userAuthDao) CountUser(ctx context.Context, condition vo.ConditionVO) (int64, error) {
	var count int64
	query := dao.db.WithContext(ctx).
		Table("tb_user_auth").
		Joins("LEFT JOIN tb_user_info ui ON tb_user_auth.user_info_id = ui.id").
		Where("1 = 1")

	if condition.Keywords != nil {
		query = query.Where("ui.nickname LIKE ?", "%"+*condition.Keywords+"%")
	}
	if condition.LoginType != nil {
		query = query.Where("tb_user_auth.login_type = ?", *condition.LoginType)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateUserAuth 更新用户认证信息
func (dao *userAuthDao) UpdateUserAuth(ctx context.Context, userAuth *model.UserAuth) error {
	return dao.db.WithContext(ctx).Save(userAuth).Error
}

// InsertUserAuth 插入用户认证信息
func (dao *userAuthDao) InsertUserAuth(ctx context.Context, userAuth *model.UserAuth) error {
	return dao.db.WithContext(ctx).Create(userAuth).Error
}

// UpdatePassword 根据用户名修改密码
func (dao *userAuthDao) UpdatePassword(ctx context.Context, username, hashedPassword string) error {
	return dao.db.WithContext(ctx).
		Model(&model.UserAuth{}).
		Where("username = ?", username).
		Update("password", hashedPassword).Error
}

// FindUserById 根据 ID 查询用户
func (dao *userAuthDao) FindUserById(ctx context.Context, id int) (*model.UserAuth, error) {
	var user model.UserAuth
	err := dao.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByUsername 根据用户名查询用户
func (dao *userAuthDao) FindUserByUsername(ctx context.Context, username string) (*model.UserAuth, error) {
	var user model.UserAuth
	err := dao.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// 记录未找到，不返回错误，返回 nil
		return nil, nil
	}
	return &user, nil
}

// FindIpSourceById 根据 ID 查询 IP 来源
func (dao *userAuthDao) FindIpSourceById(ctx context.Context, id uint) (string, error) {
	var user model.UserAuth
	err := dao.db.WithContext(ctx).Select("ip_source").First(&user, id).Error
	if err != nil {
		return "", err
	}
	return user.IPSource, nil
}

// SelectUserByUsername 根据用户名查询用户信息
func (dao *userAuthDao) SelectUserByUsername(ctx context.Context, username string) (*model.UserAuth, error) {
	var user model.UserAuth
	err := dao.db.WithContext(ctx).Select("id", "user_info_id", "username", "password", "login_type").Where(
		"username=?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateById 根据 ID 更新用户认证信息
func (dao *userAuthDao) UpdateById(ctx context.Context, userAuth *model.UserAuth) error {
	return dao.db.WithContext(ctx).Model(&model.UserAuth{}).Where("id = ?", userAuth.ID).Updates(userAuth).Error
}
