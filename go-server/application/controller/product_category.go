// Author: Steve Zhang
// Date: 2020/10/18 2:36 下午

package controller

import (
	"github.com/gin-gonic/gin"

	"go-server/application/logic"
	"go-server/common"
)

// AddProductCategoryRequest 定义新增产品类目请求结构
type AddProductCategoryRequest struct {
	ParentID       int64  `json:"parent_id" binding:""`
	CategoryName   string `json:"category_name" binding:"required,lte=128"`
	CategoryNameEN string `json:"category_name_en" binding:"required,lte=128"`
	Image          string `json:"image" binding:"lte=128"`
	Detail         string `json:"detail" binding:"lte=1024"`
	DetailEN       string `json:"detail_en" binding:"lte=1024"`
}

// AddProductCategory 新增产品类目控制器
func AddProductCategory(c *gin.Context) {
	req := &AddProductCategoryRequest{}
	if err := c.ShouldBind(req); err != nil {
		common.SetResponseContext(c, common.NewParamErrResponse(err), nil)
		return
	}
	rsp, err := logic.AddProductCategory(
		req.ParentID, req.CategoryName, req.CategoryNameEN,
		req.Image, req.Detail, req.DetailEN,
	)
	common.SetResponseContext(c, rsp, err)
	return
}

// DeleteProductCategoryRequest 定义删除产品类目请求结构
type DeleteProductCategoryRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// DeleteProductCategory 删除产品类目控制器
func DeleteProductCategory(c *gin.Context) {
	req := &DeleteProductCategoryRequest{}
	if err := c.ShouldBind(req); err != nil {
		common.SetResponseContext(c, common.NewParamErrResponse(err), nil)
		return
	}
	rsp, err := logic.DeleteProductCategory(req.ID)
	common.SetResponseContext(c, rsp, err)
	return
}

// UpdateProductCategoryRequest 定义更新产品类目请求结构
type UpdateProductCategoryRequest struct {
	ID             int64  `json:"id" binding:"required"`
	ParentID       int64  `json:"parent_id"`
	CategoryName   string `json:"category_name" binding:"required,lte=128"`
	CategoryNameEN string `json:"category_name_en" binding:"required,lte=128"`
	Image          string `json:"image" binding:"lte=128"`
	Detail         string `json:"detail" binding:"lte=1024"`
	DetailEN       string `json:"detail_en" binding:"lte=1024"`
}

// UpdateProductCategory 更新产品类目控制器
func UpdateProductCategory(c *gin.Context) {
	req := &UpdateProductCategoryRequest{}
	if err := c.ShouldBind(req); err != nil {
		common.SetResponseContext(c, common.NewParamErrResponse(err), nil)
		return
	}
	rsp, err := logic.UpdateProductCategory(
		req.ID, req.ParentID, req.CategoryName, req.CategoryNameEN,
		req.Image, req.Detail, req.DetailEN,
	)
	common.SetResponseContext(c, rsp, err)
	return
}

// QueryProductCategoryListRequest 定义查询产品类目列表请求
type QueryProductCategoryListRequest struct {
	ParentID int64 `json:"parent_id" form:"parent_id"`
}

// QueryProductCategoryList 查询产品类目控制器
func QueryProductCategoryList(c *gin.Context) {
	req := &QueryProductCategoryListRequest{}
	if err := c.ShouldBind(req); err != nil {
		common.SetResponseContext(c, common.NewParamErrResponse(err), nil)
		return
	}
	rsp, err := logic.QueryProductCategoryList(req.ParentID)
	common.SetResponseContext(c, rsp, err)
	return
}
