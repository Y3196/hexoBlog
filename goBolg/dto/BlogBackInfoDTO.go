package dto

type BlogBackInfoDTO struct {
	ViewsCount            int                    `json:"viewsCount"`
	MessageCount          int                    `json:"messageCount"`
	UserCount             int                    `json:"userCount"`
	ArticleCount          int                    `json:"articleCount"`
	CategoryDTOList       []CategoryDTO          `json:"categoryDTOList"`
	TagDTOList            []TagDTO               `json:"tagDTOList"`
	ArticleStatisticsList []ArticleStatisticsDTO `json:"articleStatisticsList"`
	UniqueViewDTOList     []UniqueViewDTO        `json:"uniqueViewDTOList"`
	ArticleRankDTOList    []ArticleRankDTO       `json:"articleRankDTOList"`
}
