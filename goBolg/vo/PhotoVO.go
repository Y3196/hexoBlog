package vo

// PhotoVO 代表照片的信息
type PhotoVO struct {
	AlbumID      int      `json:"albumId" validate:"required"`                    // 相册id
	PhotoURLList []string `json:"photoUrlList" validate:"required,dive,required"` // 照片url列表
	PhotoIDList  []int    `json:"photoIdList" validate:"required,dive,required"`  // 照片id列表
}
