package dto

// UniqueViewDTO represents the daily unique views data.
type UniqueViewDTO struct {
	Day        string `json:"day"`        // Day specifies the particular date for the view data.
	ViewsCount int    `json:"viewsCount"` // ViewsCount is the number of views recorded on that day.
}
