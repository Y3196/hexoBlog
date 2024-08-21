package dto

// ArticleStatisticsDTO contains statistics of articles over a specific date.
type ArticleStatisticsDTO struct {
	Date  string `json:"date"`  // Date represents the date of the statistics
	Count int    `json:"count"` // Count is the number of occurrences or articles on that date
}
