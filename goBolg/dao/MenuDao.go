package dao

import (
	"context"
	"goBolg/model"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
)

// MenuDao 菜单 DAO 接口
type MenuDao interface {
	// 根据用户 ID 查询菜单
	ListMenusByUserInfoID(ctx context.Context, userInfoID uint) ([]model.Menu, error)

	// 查询菜单列表（可选条件）
	ListMenus(ctx context.Context, condition vo.ConditionVO) ([]model.Menu, error)

	// 查询子菜单并删除
	DeleteSubMenus(ctx context.Context, menuID int) error

	FindById(ctx context.Context, id uint) (*model.Menu, error)

	Save(ctx context.Context, menu *model.Menu) error

	Update(ctx context.Context, menu *model.Menu) error
}

type menuDao struct {
	db *gorm.DB
}

// NewMenuDao 创建新的 MenuDao 实例
func NewMenuDao(db *gorm.DB) MenuDao {
	return &menuDao{db: db}
}

// ListMenusByUserInfoID 根据用户 ID 查询菜单
func (dao *menuDao) ListMenusByUserInfoID(ctx context.Context, userInfoID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := dao.db.WithContext(ctx).
		Raw(`
			SELECT DISTINCT
				m.id,
				m.name,
				m.path,
				m.component,
				m.icon,
				m.is_hidden,
				m.parent_id,
				m.order_num
			FROM
				tb_user_role ur
				JOIN tb_role_menu rm ON ur.role_id = rm.role_id
				JOIN tb_menu m ON rm.menu_id = m.id
			WHERE
				ur.user_id = ?`, userInfoID).
		Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// ListMenus 查询菜单列表（可选条件）
func (dao *menuDao) ListMenus(ctx context.Context, condition vo.ConditionVO) ([]model.Menu, error) {
	var menus []model.Menu
	query := dao.db.WithContext(ctx).Model(&model.Menu{})

	if condition.Keywords != nil && *condition.Keywords != "" {
		query = query.Where("name LIKE ?", "%"+*condition.Keywords+"%")
	}

	err := query.Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// DeleteSubMenus 查询子菜单并删除
func (dao *menuDao) DeleteSubMenus(ctx context.Context, menuID int) error {
	var subMenuIDs []int
	// 查询子菜单 ID
	err := dao.db.WithContext(ctx).Model(&model.Menu{}).
		Select("id").
		Where("parent_id = ?", menuID).
		Pluck("id", &subMenuIDs).Error
	if err != nil {
		return err
	}
	// 添加当前菜单 ID
	subMenuIDs = append(subMenuIDs, menuID)
	// 删除子菜单
	err = dao.db.WithContext(ctx).Delete(&model.Menu{}, subMenuIDs).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *menuDao) Save(ctx context.Context, menu *model.Menu) error {
	log.Printf("Creating new menu: %+v", menu)
	return dao.db.WithContext(ctx).Create(menu).Error
}

func (dao *menuDao) Update(ctx context.Context, menu *model.Menu) error {
	log.Printf("Updating existing menu: %+v", menu)
	return dao.db.WithContext(ctx).Save(menu).Error
}

func (dao *menuDao) FindById(ctx context.Context, id uint) (*model.Menu, error) {
	var menu model.Menu
	err := dao.db.WithContext(ctx).First(&menu, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("No record found for ID: %d", id)
			return nil, err
		}
		log.Printf("Error finding menu by ID: %v", err)
		return nil, err
	}
	log.Printf("Record found: %+v", menu)
	return &menu, nil
}
