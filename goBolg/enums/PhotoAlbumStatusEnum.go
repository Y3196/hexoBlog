package enums

// PhotoAlbumStatusEnum represents the status of a photo album
type PhotoAlbumStatusEnum struct {
	Status int
	Desc   string
}

var (
	// PUBLIC_STATUS represents a public photo album
	PUBLIC_STATUS = PhotoAlbumStatusEnum{1, "公开"}
	// SECRET_STATUS represents a private photo album
	SECRET_STATUS = PhotoAlbumStatusEnum{2, "私密"}
)

// GetPhotoAlbumStatusEnums returns a list of all photo album status enums
func GetPhotoAlbumStatusEnums() []PhotoAlbumStatusEnum {
	return []PhotoAlbumStatusEnum{
		PUBLIC_STATUS,
		SECRET_STATUS,
	}
}
