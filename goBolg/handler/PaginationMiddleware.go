package handler

import (
	"github.com/gin-gonic/gin"
	"goBolg/utils"
	"strconv"
)

func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		size := c.DefaultQuery("size", "10")

		currentPage, err := strconv.Atoi(page)
		if err != nil || currentPage < 1 {
			currentPage = 1
		}

		pageSize, err := strconv.Atoi(size)
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		ctx := utils.SetCurrentPage(c.Request.Context(), &utils.Page{Current: currentPage, Size: pageSize})
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
