package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-server/application/controller"
	"go-server/application/middleware"
)

// setupRouter 设置路由
func setupRouter() http.Handler {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.Recovery)

	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		_, _ = c.Writer.WriteString("ok")
		c.AbortWithStatus(http.StatusOK)
	})

	// 接口路由分组
	api := router.Group("/api").Use(
		middleware.Log,
		middleware.Response,
	)

	// 以下接口需要权限认证
	api.Use(middleware.Auth)
	{
		api.POST("/product.category.add", controller.AddProductCategory)       // 新增产品类目
		api.POST("/product.category.delete", controller.DeleteProductCategory) // 删除产品类目
		api.POST("/product.category.update", controller.UpdateProductCategory) // 更新产品类目
		api.GET("/product.category.list", controller.QueryProductCategoryList) // 查询产品类目列表
	}

	return router
}
