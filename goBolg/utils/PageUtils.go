package utils

import (
	"context"
)

type Page struct {
	Current int
	Size    int
}

type pageKey struct{}

var pageHolder = &Page{Current: 1, Size: 10} // 默认分页值

// SetCurrentPage 设置当前分页
func SetCurrentPage(ctx context.Context, page *Page) context.Context {
	return context.WithValue(ctx, pageKey{}, page)
}

// GetPage 获取当前分页
func GetPage(ctx context.Context) *Page {
	if page, ok := ctx.Value(pageKey{}).(*Page); ok {
		return page
	}
	return pageHolder
}

// GetCurrent 获取当前页码
func GetCurrent(ctx context.Context) int {
	return int(GetPage(ctx).Current)
}

// GetSize 获取每页大小
func GetSize(ctx context.Context) int {
	return int(GetPage(ctx).Size)
}

// GetLimitCurrent 获取分页偏移量
func GetLimitCurrent(ctx context.Context) int {
	return (GetCurrent(ctx) - 1) * GetSize(ctx)
}
