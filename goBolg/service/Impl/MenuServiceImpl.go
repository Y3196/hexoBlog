package Impl

import (
	"context"
	"errors"
	"goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
	"sort"
	"time"
)

// menuServiceImpl 实现
type menuServiceImpl struct {
	menuDao     dao.MenuDao
	roleMenuDao dao.RoleMenuDao
}

// NewMenuService 创建新的 MenuService
func NewMenuService(menuDao dao.MenuDao, roleMenuDao dao.RoleMenuDao) service.MenuService {
	return &menuServiceImpl{menuDao: menuDao, roleMenuDao: roleMenuDao}
}

// ListMenus 查看菜单列表
func (s *menuServiceImpl) ListMenus(ctx context.Context, conditionVO vo.ConditionVO) ([]dto.MenuDTO, error) {
	// 查询菜单数据
	menuList, err := s.menuDao.ListMenus(ctx, conditionVO)
	if err != nil {
		return nil, err
	}

	// 获取目录列表
	catalogList := listCatalog(menuList)

	// 获取目录下的子菜单
	childrenMap := getMenuMap(menuList)

	// 组装目录菜单数据
	var menuDTOList []dto.MenuDTO
	for _, item := range catalogList {
		var menuDTO dto.MenuDTO
		utils.BeanCopyObject(item, &menuDTO)

		// 获取目录下的菜单排序
		children := childrenMap[item.ID]
		sort.Slice(children, func(i, j int) bool {
			return children[i].OrderNum < children[j].OrderNum
		})
		childrenDTOList := utils.BeanCopyList(children, &dto.MenuDTO{}).([]dto.MenuDTO)

		menuDTO.Children = childrenDTOList
		delete(childrenMap, item.ID)
		menuDTOList = append(menuDTOList, menuDTO)
	}

	// 若还有菜单未取出则拼接
	if len(childrenMap) > 0 {
		var childrenList []model.Menu
		for _, children := range childrenMap {
			childrenList = append(childrenList, children...)
		}
		sort.Slice(childrenList, func(i, j int) bool {
			return childrenList[i].OrderNum < childrenList[j].OrderNum
		})
		childrenDTOList := utils.BeanCopyList(childrenList, &dto.MenuDTO{}).([]dto.MenuDTO)
		menuDTOList = append(menuDTOList, childrenDTOList...)
	}

	return menuDTOList, nil
}

// listCatalog 获取目录列表
func listCatalog(menuList []model.Menu) []model.Menu {
	var catalogList []model.Menu
	for _, item := range menuList {
		if item.ParentID == nil {
			catalogList = append(catalogList, item)
		}
	}

	sort.Slice(catalogList, func(i, j int) bool {
		return catalogList[i].OrderNum < catalogList[j].OrderNum
	})

	return catalogList
}

// SaveOrUpdateMenu 保存或更新菜单
func (s *menuServiceImpl) SaveOrUpdateMenu(ctx context.Context, menuVO vo.MenuVO) error {
	menu := &model.Menu{}
	utils.BeanCopy(menuVO, menu)

	log.Printf("Saving or updating menu: %+v", menu)

	if menu.ID != 0 {
		existingMenu, err := s.menuDao.FindById(ctx, menu.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				menu.CreateTime = time.Now()
				log.Printf("Record not found, creating new menu with ID: %d", menu.ID)
				return s.menuDao.Save(ctx, menu)
			} else {
				log.Printf("Error finding existing menu: %v", err)
				return err
			}
		} else {
			menu.CreateTime = existingMenu.CreateTime
			now := time.Now()
			menu.UpdateTime = &now
			log.Printf("Record found, updating existing menu with ID: %d", menu.ID)
			return s.menuDao.Update(ctx, menu)
		}
	} else {
		menu.CreateTime = time.Now()
		log.Printf("Creating new menu with ID: %d", menu.ID)
		return s.menuDao.Save(ctx, menu)
	}
}

// DeleteMenu 删除菜单
func (s *menuServiceImpl) DeleteMenu(ctx context.Context, menuId uint) error {
	// 查询是否有角色关联
	count, err := s.roleMenuDao.CountByMenuID(ctx, menuId)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("菜单下有角色关联")
	}

	// 查询子菜单并删除
	err = s.menuDao.DeleteSubMenus(ctx, int(menuId))
	if err != nil {
		return err
	}

	log.Printf("Deleted menu with ID: %d", menuId)
	return nil
}

func (s *menuServiceImpl) ListMenuOptions(ctx context.Context) ([]dto.LabelOptionDTO, error) {
	condition := vo.ConditionVO{}
	menuList, err := s.menuDao.ListMenus(ctx, condition)
	if err != nil {
		return nil, err
	}
	catalogList := listCatalog(menuList)
	childrenMap := getMenuMap(menuList)

	var menuOptions []dto.LabelOptionDTO
	for _, item := range catalogList {
		var childrenOptions []dto.LabelOptionDTO
		if children, ok := childrenMap[item.ID]; ok {
			for _, child := range children {
				childrenOptions = append(childrenOptions, dto.LabelOptionDTO{
					ID:    child.ID,
					Label: child.Name,
				})
			}
		}
		menuOptions = append(menuOptions, dto.LabelOptionDTO{
			ID:       item.ID,
			Label:    item.Name,
			Children: childrenOptions,
		})
	}
	return menuOptions, nil
}

// ListUserMenus 查询用户菜单列表
func (s *menuServiceImpl) ListUserMenus(ctx context.Context, userID uint) ([]dto.UserMenuDTO, error) {
	menuList, err := s.menuDao.ListMenusByUserInfoID(ctx, uint(userID))
	if err != nil {
		return nil, err
	}
	catalogList := listCatalog(menuList)
	childrenMap := getMenuMap(menuList)

	return s.convertUserMenuList(catalogList, childrenMap), nil
}

func (s *menuServiceImpl) convertUserMenuList(catalogList []model.Menu, childrenMap map[uint][]model.Menu) []dto.UserMenuDTO {
	var userMenuDTOs []dto.UserMenuDTO
	for _, item := range catalogList {
		userMenuDTO := dto.UserMenuDTO{}
		utils.BeanCopyObject(item, &userMenuDTO)

		if item.IsHidden != nil {
			userMenuDTO.Hidden = *item.IsHidden == constants.True
		} else {
			userMenuDTO.Hidden = false
		}

		var list []dto.UserMenuDTO
		children := childrenMap[item.ID]
		for _, menu := range children {
			dto := dto.UserMenuDTO{}
			utils.BeanCopyObject(menu, &dto)
			if menu.IsHidden != nil {
				dto.Hidden = *menu.IsHidden == constants.True
			} else {
				dto.Hidden = false
			}
			list = append(list, dto)
		}

		if len(list) == 0 {
			userMenuDTO.Path = item.Path
			userMenuDTO.Component = constants.Component
			list = append(list, dto.UserMenuDTO{
				Path:      "",
				Name:      item.Name,
				Icon:      item.Icon,
				Component: item.Component,
			})
		}

		userMenuDTO.Children = list
		userMenuDTOs = append(userMenuDTOs, userMenuDTO)
	}
	return userMenuDTOs
}

// getMenuMap 获取目录下的子菜单
func getMenuMap(menuList []model.Menu) map[uint][]model.Menu {
	menuMap := make(map[uint][]model.Menu)
	for _, menu := range menuList {
		if menu.ParentID != nil {
			parentID := *menu.ParentID
			menuMap[parentID] = append(menuMap[parentID], menu)
		}
	}
	return menuMap
}
