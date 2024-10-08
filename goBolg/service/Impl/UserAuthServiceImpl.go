package Impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"goBolg/constant"
	"goBolg/dao"
	"goBolg/dto"
	"goBolg/enums"
	"goBolg/model"
	"goBolg/rabbitmq/rabbitService"
	"goBolg/service"
	"goBolg/utils"
	"goBolg/vo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

// userAuthServiceImpl 是 UserAuthService 的实现
type userAuthServiceImpl struct {
	emailService    rabbitService.EmailService
	redisService    service.RedisService
	rabbitService   rabbitService.RabbitService
	userInfoDao     dao.UserInfoDao
	userRoleDao     dao.UserRoleDao
	userAuthDao     dao.UserAuthDao
	blogInfoService service.BlogInfoService
}

// NewUserAuthService 创建新的 UserAuthService 实例
func NewUserAuthService(emailService rabbitService.EmailService, redisService service.RedisService, rabbitService rabbitService.RabbitService, userInfoDao dao.UserInfoDao, userRoleDao dao.UserRoleDao, userAuthDao dao.UserAuthDao, blogInfoService service.BlogInfoService) service.UserAuthService {
	return &userAuthServiceImpl{
		emailService:    emailService,
		redisService:    redisService,
		rabbitService:   rabbitService,
		userInfoDao:     userInfoDao,
		userRoleDao:     userRoleDao,
		userAuthDao:     userAuthDao,
		blogInfoService: blogInfoService,
	}
}

// SendCode 发送验证码
func (s *userAuthServiceImpl) SendCode(username string) error {
	// 校验邮箱合法性
	if !utils.CheckEmail(username) {
		return fmt.Errorf("请输入正确邮箱")
	}

	// 生成六位随机验证码
	code := utils.GetRandomCode()

	// 创建邮件内容
	subject := "验证码"
	content := fmt.Sprintf("您的验证码为 %s 有效期15分钟，请不要告诉他人哦！", code)

	// 发送邮件
	err := s.emailService.SendEmail(username, subject, content)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	// 将验证码存入 Redis，设置过期时间为15分钟
	err = s.redisService.Set(context.Background(), constants.UserCodeKey+username, code, 15*time.Minute)
	if err != nil {
		log.Printf("Failed to store code in Redis: %v", err)
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	return nil
}

// Register registers a new user.
func (s *userAuthServiceImpl) Register(ctx context.Context, user vo.UserVO) error {
	// 校验账号是否合法
	exists, err := s.CheckUser(ctx, user)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("邮箱已被注册！")
	}

	// 获取用户头像
	websiteConfig, err := s.blogInfoService.GetWebsiteConfig(ctx)
	if err != nil {
		return err
	}

	// 新增用户信息
	userInfo := &model.UserInfo{
		Email:    user.Username,
		Nickname: constants.DefaultNickname + generateUniqueID(),
		Avatar:   websiteConfig.UserAvatar,
	}
	if err := s.userInfoDao.InsertUserInfo(ctx, userInfo); err != nil {
		return err
	}

	// 绑定用户角色
	userRole := &model.UserRole{
		UserID: userInfo.ID,
		RoleID: enums.USER.RoleID,
	}
	if err := s.userRoleDao.InsertUserRole(ctx, userRole); err != nil {
		return err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 新增用户账号
	userAuth := &model.UserAuth{
		UserInfoID: userInfo.ID,
		Username:   user.Username,
		Password:   string(hashedPassword),
		LoginType:  enums.EMAIL.Type,
	}
	return s.userAuthDao.InsertUserAuth(ctx, userAuth)
}

func (s *userAuthServiceImpl) CheckUser(ctx context.Context, user vo.UserVO) (bool, error) {
	// 验证验证码
	storedCode, err := s.redisService.Get(ctx, constants.UserCodeKey+user.Username)
	if err != nil {
		return false, fmt.Errorf("failed to get code from redis: %v", err)
	}

	if user.Code != storedCode {
		return false, fmt.Errorf("验证码错误")
	}

	// 查询用户名是否存在
	userAuth, err := s.userAuthDao.FindUserByUsername(ctx, user.Username)
	if err != nil {
		return false, fmt.Errorf("failed to query user: %v", err)
	}

	return userAuth != nil, nil
}

// generateUniqueID 生成唯一ID（伪代码，根据实际情况实现）
func generateUniqueID() string {
	// 这里需要一个函数来生成唯一ID，你可以使用UUID，或者根据你的实际需求实现
	return "unique-id"
}

// UpdatePassword 更新用户密码
func (s *userAuthServiceImpl) UpdatePassword(ctx context.Context, user *vo.UserVO) error {
	// 校验账号是否合法
	existingUser, err := s.userAuthDao.FindUserByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("邮箱尚未注册！")
	}

	// 根据用户名修改密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userAuthDao.UpdatePassword(ctx, user.Username, string(hashedPassword))
}

// ListUserAreas fetches the user area distribution based on the provided condition.
func (s *userAuthServiceImpl) ListUserAreas(ctx context.Context, conditionVO vo.ConditionVO) ([]dto.UserAreaDTO, error) {
	var userAreaDTOList []dto.UserAreaDTO

	areaType := enums.GetUserAreaType(*conditionVO.Type)
	if areaType == nil {
		return userAreaDTOList, nil // or return an error if needed
	}

	switch *areaType {
	case enums.USERTYPE:
		// 查询注册用户区域分布
		userArea, err := s.redisService.Get(ctx, constants.UserArea)
		if err != nil && err != redis.Nil {
			return nil, err
		}
		log.Printf("Raw data from Redis: %s", userArea)

		if userArea != "" {
			// 使用 DecodeDoubleEscapedJSON 处理双重转义的 JSON 字符串
			decodedJSON, err := utils.DecodeDoubleEscapedJSON(userArea)
			if err != nil {
				log.Printf("Error decoding double escaped JSON: %v", err)
				return nil, err
			}

			// 将解码后的字符串反序列化为结构体数组
			err = json.Unmarshal([]byte(decodedJSON), &userAreaDTOList)
			if err != nil {
				log.Printf("Error unmarshaling JSON into struct: %v", err)
				return nil, err
			}
		}

	case enums.VISITOR:
		// 查询游客区域分布
		visitorArea, err := s.redisService.HGetAll(ctx, constants.VisitorArea)
		if err != nil && err != redis.Nil {
			return nil, err
		}
		if len(visitorArea) > 0 {
			for name, value := range visitorArea {
				val, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, err
				}
				userAreaDTOList = append(userAreaDTOList, dto.UserAreaDTO{
					Name:  name,
					Value: val,
				})
			}
		}
	default:
		// 默认情况，返回空的 userAreaDTOList
	}

	return userAreaDTOList, nil
}
