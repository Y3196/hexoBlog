package Impl

import (
	"context"
	"errors"
	"fmt"
	"goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"log"
	"strconv"
)

// talkService 实现 TalkService 接口
type talkServiceImpl struct {
	talkDao      dao.TalkDao
	commentDao   dao.CommentDao
	redisService service.RedisService
}

// NewTalkService 创建新的 TalkService 实例
func NewTalkService(talkDao dao.TalkDao, commentDao dao.CommentDao, redisService service.RedisService) service.TalkService {
	return &talkServiceImpl{talkDao: talkDao, commentDao: commentDao, redisService: redisService}
}

// ListHomeTalks 获取首页说说
func (s *talkServiceImpl) ListHomeTalks(ctx context.Context) ([]string, error) {
	return s.talkDao.ListLatestTalks(ctx)
}

// GetTalkById 根据ID获取说说信息
// GetTalkById 根据ID获取说说信息
func (s *talkServiceImpl) GetTalkById(ctx context.Context, talkId int) (*dto.TalkDTO, error) {
	talkDTO, err := s.talkDao.GetTalkById(ctx, talkId)
	if err != nil {
		return nil, err
	}
	if talkDTO == nil {
		return nil, errors.New("说说不存在")
	}

	likeCount, err := s.redisService.HGet(ctx, constants.TalkLikeCount, strconv.Itoa(talkId))
	if err != nil {
		talkDTO.LikeCount = 0 // 默认为0，避免解析错误
	} else {
		talkDTO.LikeCount, _ = strconv.Atoi(likeCount) // 忽略错误，因为已经做了默认值处理
	}
	// 转换图片格式
	if talkDTO.Images != "" {
		imgList, err := utils.CastList(talkDTO.Images, utils.StringConstructor)
		if err == nil {
			talkDTO.ImgList = imgList
		}
	}

	return talkDTO, nil
}

// ListTalks 查询所有说说，并统计评论量和点赞量
func (s *talkServiceImpl) ListTalks(ctx context.Context) (vo.PageResult[dto.TalkDTO], error) {
	// 查询说说总量
	count, err := s.talkDao.CountAll(ctx)
	if err != nil {
		return vo.PageResult[dto.TalkDTO]{}, err
	}
	if count == 0 {
		return vo.PageResult[dto.TalkDTO]{}, nil
	}

	// 分页查询说说
	talkDTOList, err := s.talkDao.ListTalks(ctx, utils.GetLimitCurrent(ctx), utils.GetSize(ctx))
	if err != nil {
		return vo.PageResult[dto.TalkDTO]{}, err
	}

	// 查询说说评论量
	talkIdList := make([]int, len(talkDTOList))
	for i, talk := range talkDTOList {
		talkIdList[i] = talk.ID
	}
	commentCounts, err := s.commentDao.ListCommentCountByTopicIds(ctx, talkIdList)
	if err != nil {
		return vo.PageResult[dto.TalkDTO]{}, err
	}
	commentCountMap := convertCommentCountListToMap(commentCounts)

	// 查询说说点赞量
	likeCountMap, err := s.redisService.HGetAll(ctx, constants.TalkLikeCount)
	if err != nil {
		return vo.PageResult[dto.TalkDTO]{}, err
	}

	for _, item := range talkDTOList {
		if likeCount, ok := likeCountMap[strconv.Itoa(item.ID)]; ok {
			if count, err := strconv.Atoi(likeCount); err == nil {
				item.LikeCount = count
			}
		}
		item.CommentCount = commentCountMap[item.ID] // 确保 commentCountMap 的值是 int 类型
		if item.Images != "" {
			imgList, err := utils.CastList(item.Images, utils.StringConstructor)
			if err == nil {
				item.ImgList = imgList
			}
		}
	}

	return vo.PageResult[dto.TalkDTO]{RecordList: talkDTOList, Count: count}, nil
}

// convertCommentCountListToMap 将评论数量列表转换为映射
func convertCommentCountListToMap(commentCounts []dto.CommentCountDTO) map[int]int {
	countMap := make(map[int]int)
	for _, cc := range commentCounts {
		countMap[cc.ID] = cc.CommentCount
	}
	return countMap
}

func (s *talkServiceImpl) SaveTalkLike(ctx context.Context, talkId int) (int, error) {
	user, ok := utils.GetLoginUser(ctx)
	if !ok {
		return 0, fmt.Errorf("用户未登录")
	}

	talkLikeKey := constants.TalkUserLike + strconv.Itoa(user.UserInfoID)
	talkIdStr := strconv.Itoa(talkId)

	isMember, err := s.redisService.SIsMember(ctx, talkLikeKey, talkIdStr)
	if err != nil {
		return 0, err
	}

	if isMember {
		_, err := s.redisService.SRemove(ctx, talkLikeKey, talkIdStr)
		if err != nil {
			return 0, err
		}
		_, err = s.redisService.HDecr(ctx, constants.TalkLikeCount, talkIdStr, 1)
		if err != nil {
			return 0, err
		}
		log.Printf("Successfully removed talk ID %d from set for user %d", talkId, user.UserInfoID)
	} else {
		_, err := s.redisService.SAdd(ctx, talkLikeKey, talkIdStr)
		if err != nil {
			log.Printf("Error adding talk ID to set for user %d: %v", user.UserInfoID, err)
			return 0, err
		}
		log.Printf("Successfully added talk ID %d to set for user %d", talkId, user.UserInfoID)
		_, err = s.redisService.HIncr(ctx, constants.TalkLikeCount, talkIdStr, 1)
		if err != nil {
			return 0, err
		}
	}

	likeCountStr, err := s.redisService.HGet(ctx, constants.TalkLikeCount, talkIdStr)
	if err != nil {
		return 0, err
	}
	likeCount, _ := strconv.Atoi(likeCountStr) // 这里忽略转换错误，因为前面应该已经处理过了
	return likeCount, nil
}
