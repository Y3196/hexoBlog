package dto

// LabelOptionDTO 代表标签选项 DTO
type LabelOptionDTO struct {
	ID       uint             `json:"id"`       // 选项 ID
	Label    string           `json:"label"`    // 选项名
	Children []LabelOptionDTO `json:"children"` // 子选项
}
