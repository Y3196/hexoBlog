package Impl

import (
	"context"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/model"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"net/http"
)

// messageServiceImpl implements MessageService
type messageServiceImpl struct {
	messageDao      dao.MessageDao
	blogInfoService service.BlogInfoService
	request         *http.Request
}

// NewMessageService creates a new MessageService
func NewMessageService(messageDao dao.MessageDao, blogInfoService service.BlogInfoService, request *http.Request) service.MessageService {
	return &messageServiceImpl{
		messageDao:      messageDao,
		blogInfoService: blogInfoService,
		request:         request,
	}
}

// SaveMessage saves a message
func (s *messageServiceImpl) SaveMessage(ctx context.Context, messageVO vo.MessageVO) error {
	websiteConfig, err := s.blogInfoService.GetWebsiteConfig(ctx)
	if err != nil {
		return err
	}

	// 判断是否需要审核
	var isReview int = websiteConfig.IsMessageReview

	// 获取用户ip
	ipAddress := utils.GetIPAddress(s.request)
	ipSource := utils.GetIPSource(ipAddress)
	message := model.Message{}
	utils.BeanCopyObject(messageVO, &message)
	message.MessageContent = utils.HTMLFilter(message.MessageContent)
	message.IPAddress = ipAddress
	message.IsReview = isReview
	message.IPSource = ipSource
	return s.messageDao.Insert(ctx, message)
}

// ListMessages 查询审核通过的留言列表
func (s *messageServiceImpl) ListMessages(ctx context.Context) ([]dto.MessageDTO, error) {
	// 查询留言列表
	messageList, err := s.messageDao.ListMessages(ctx, true)
	if err != nil {
		return nil, err
	}
	return utils.BeanCopyList(messageList, &dto.MessageDTO{}).([]dto.MessageDTO), nil
}

// ListMessageBackDTO 分页查询留言列表
func (s *messageServiceImpl) ListMessageBackDTO(ctx context.Context, condition vo.ConditionVO) (vo.PageResult[dto.MessageBackDTO], error) {
	// 获取当前页码和每页大小
	currentPage := utils.GetCurrent(ctx)
	pageSize := utils.GetSize(ctx)

	// 创建分页查询条件
	messageList, totalRecords, err := s.messageDao.PageMessages(ctx, currentPage, pageSize, condition)
	if err != nil {
		return vo.PageResult[dto.MessageBackDTO]{}, err
	}

	// 转换DTO
	var messageBackDTOList []dto.MessageBackDTO
	for _, message := range messageList {
		var messageBackDTO dto.MessageBackDTO
		utils.BeanCopy(message, &messageBackDTO)
		messageBackDTOList = append(messageBackDTOList, messageBackDTO)
	}

	// 返回分页结果
	return vo.NewPageResult(messageBackDTOList, int(totalRecords)), nil
}

// UpdateMessagesReview 更新留言审核状态
func (s *messageServiceImpl) UpdateMessagesReview(ctx context.Context, reviewVO vo.ReviewVO) error {
	var messageList []model.Message
	for _, id := range reviewVO.IDList {
		messageList = append(messageList, model.Message{
			ID:       id,
			IsReview: reviewVO.IsReview,
		})
	}
	return s.messageDao.UpdateBatchById(ctx, messageList)
}

// DeleteMessages 删除留言
func (s *messageServiceImpl) DeleteMessages(ctx context.Context, messageIdList []uint) error {
	return s.messageDao.RemoveByIds(ctx, messageIdList)
}
