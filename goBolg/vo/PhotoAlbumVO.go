package vo

// PhotoAlbumVO 代表相册的信息
type PhotoAlbumVO struct {
	ID         int    `json:"id"`                             // 相册id
	AlbumName  string `json:"albumName" validate:"required"`  // 相册名
	AlbumDesc  string `json:"albumDesc"`                      // 相册描述
	AlbumCover string `json:"albumCover" validate:"required"` // 相册封面
	Status     int    `json:"status"`                         // 状态值 1公开 2私密
}
