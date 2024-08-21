package dto

// PhotoAlbumDTO 代表相册 DTO
type PhotoAlbumDTO struct {
	ID         int    `json:"id"`         // 相册id
	AlbumName  string `json:"albumName"`  // 相册名
	AlbumDesc  string `json:"albumDesc"`  // 相册描述
	AlbumCover string `json:"albumCover"` // 相册封面
}
