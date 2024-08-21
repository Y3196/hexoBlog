package vo

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"strconv"
)

// CommentVO 代表评论信息
type CommentVO struct {
	ReplyUserID    int    `json:"replyUserId,omitempty"`              // 回复用户id
	TopicID        int    `json:"topicId,omitempty"`                  // 评论主题id
	CommentContent string `json:"commentContent" validate:"required"` // 评论内容，必填
	ParentID       int    `json:"parentId,omitempty"`                 // 父评论id
	Type           int    `json:"type" validate:"required"`           // 类型，必填
}

// ValidateCommentVO 验证 CommentVO 实例
func ValidateCommentVO(comment CommentVO) error {
	validate := validator.New()
	return validate.Struct(comment)
}

// Custom Unmarshaler to handle string to int conversion for TopicID
func (c *CommentVO) UnmarshalJSON(data []byte) error {
	var aux struct {
		ReplyUserID    int    `json:"replyUserId,omitempty"`
		TopicID        string `json:"topicId,omitempty"`
		CommentContent string `json:"commentContent" validate:"required"`
		ParentID       int    `json:"parentId,omitempty"`
		Type           int    `json:"type" validate:"required"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.ReplyUserID = aux.ReplyUserID
	c.CommentContent = aux.CommentContent
	c.ParentID = aux.ParentID
	c.Type = aux.Type

	// Convert TopicID from string to int
	if aux.TopicID != "" {
		topicID, err := strconv.Atoi(aux.TopicID)
		if err != nil {
			return err
		}
		c.TopicID = topicID
	}

	return nil
}
