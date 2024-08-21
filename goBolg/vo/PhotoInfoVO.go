package vo

// PhotoInfoVO 代表照片的信息
type PhotoInfoVO struct {
	ID        int    `json:"id" validate:"required"`        // 照片id
	PhotoName string `json:"photoName" validate:"required"` // 照片名
	PhotoDesc string `json:"photoDesc"`                     // 照片描述
}
