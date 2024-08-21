package Impl

import (
	"context"
	"errors"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"gorm.io/gorm"
	"log"
	"time"
)

type friendLinkServiceImpl struct {
	friendLinkDao dao.FriendLinkDao
}

func NewFriendLinkServiceImpl(friendLinkDao dao.FriendLinkDao) service.FriendLinkService {
	return &friendLinkServiceImpl{friendLinkDao: friendLinkDao}
}

func (s *friendLinkServiceImpl) ListFriendLinks(ctx context.Context) ([]dto.FriendLinkDTO, error) {
	friendLinkList, err := s.friendLinkDao.SelectList(ctx)
	if err != nil {
		return nil, err
	}

	var friendLinkDTOs []dto.FriendLinkDTO
	for _, link := range friendLinkList {
		friendLinkDTOs = append(friendLinkDTOs, dto.FriendLinkDTO{
			ID:          link.ID,
			LinkName:    link.LinkName,
			LinkAddress: link.LinkAddress,
			LinkAvatar:  link.LinkAvatar,
			LinkIntro:   link.LinkIntro,
		})
	}

	return friendLinkDTOs, nil
}

func (s *friendLinkServiceImpl) ListFriendLinkDTO(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.FriendLinkBackDTO], error) {
	current := utils.GetCurrent(ctx)
	size := utils.GetSize(ctx)
	friendLinks, count, err := s.friendLinkDao.PagedFriendLinks(ctx, current, size, condition)
	if err != nil {
		log.Printf("Error retrieving friend links: %v", err)
		return vo.PageResult[dto.FriendLinkBackDTO]{}, err
	}

	friendLinkBackDTOList := utils.BeanCopyList(friendLinks, []dto.FriendLinkBackDTO{}).([]dto.FriendLinkBackDTO)
	return vo.NewPageResult(friendLinkBackDTOList, int(count)), nil
}

// SaveOrUpdateFriendLink 保存或更新友链
func (s *friendLinkServiceImpl) SaveOrUpdateFriendLink(ctx context.Context, friendLinkVO vo.FriendLinkVO) error {
	friendLink := &model.FriendLink{
		ID:          friendLinkVO.ID,
		LinkName:    friendLinkVO.LinkName,
		LinkAvatar:  friendLinkVO.LinkAvatar,
		LinkAddress: friendLinkVO.LinkAddress,
		LinkIntro:   friendLinkVO.LinkIntro,
	}

	log.Printf("Received FriendLinkVO: %+v", friendLinkVO)
	log.Printf("Saving or updating friend link: %+v", friendLink)

	if friendLink.ID != 0 {
		existing, err := s.friendLinkDao.FindById(ctx, friendLink.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("Record not found, creating new friend link")
				friendLink.CreateTime = time.Now()
				return s.friendLinkDao.Save(ctx, friendLink)
			}
			log.Printf("Error finding friend link by ID: %v", err)
			return err
		}
		log.Printf("Existing friend link found: %+v", existing)
		friendLink.CreateTime = existing.CreateTime
		return s.friendLinkDao.Update(ctx, friendLink)
	} else {
		log.Println("Creating new friend link")
		friendLink.CreateTime = time.Now()
		return s.friendLinkDao.Save(ctx, friendLink)
	}
}

func (s *friendLinkServiceImpl) RemoveFriendLinks(ctx context.Context, ids []uint) error {
	return s.friendLinkDao.RemoveByIds(ctx, ids)
}
