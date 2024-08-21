package dto

// PhotoAlbumBackDTO 代表后台相册 DTO
type PhotoAlbumBackDTO struct {
	ID         int    `json:"id"`         // 相册id
	AlbumName  string `json:"albumName"`  // 相册名
	AlbumDesc  string `json:"albumDesc"`  // 相册描述
	AlbumCover string `json:"albumCover"` // 相册封面
	PhotoCount int    `json:"photoCount"` // 照片数量
	Status     int    `json:"status"`     // 状态值 1公开 2私密
}
