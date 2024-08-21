package dto

// PhotoDTO 代表照片 DTO
type PhotoDTO struct {
	PhotoAlbumCover string   `json:"photoAlbumCover"` // 相册封面
	PhotoAlbumName  string   `json:"photoAlbumName"`  // 相册名
	PhotoList       []string `json:"photoList"`       // 照片列表
}
