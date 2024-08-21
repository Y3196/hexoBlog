package dto

// PhotoBackDTO 代表后台照片 DTO
type PhotoBackDTO struct {
	ID        int    `json:"id"`        // 照片id
	PhotoName string `json:"photoName"` // 照片名
	PhotoDesc string `json:"photoDesc"` // 照片描述
	PhotoSrc  string `json:"photoSrc"`  // 照片地址
}
