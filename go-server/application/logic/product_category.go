// Author: Steve Zhang
// Date: 2020/10/18 2:21 下午

package logic

import (
	"fmt"

	"go-server/common"
	"go-server/model"
)

// AddProductCategory 新增产品类目逻辑
func AddProductCategory(parentID int64, categoryName, categoryNameEN, image, desc, descEN string) (rsp *common.Response, err error) {
	rsp = common.NewOKResponse()

	cate := &model.ProductCategory{
		ParentID:       parentID,
		CategoryName:   categoryName,
		CategoryNameEN: categoryNameEN,
		Image:          image,
		Detail:         desc,
		DetailEN:       descEN,
	}

	id, err := model.AddProductCategory(cate)
	if err != nil {
		rsp.Code = common.ResponseCodeInternalErr
		rsp.Message = "新增产品类目失败"
		err = fmt.Errorf("model.AddProductCategory[category=%+v]: %w", *cate, err)
		return
	}

	rsp.Data = common.ResponseData{
		"category_id": id,
	}

	return
}

// DeleteProductCategory 删除产品类目逻辑
func DeleteProductCategory(id int64) (rsp *common.Response, err error) {
	rsp = common.NewOKResponse()

	if err = model.DeleteProductCategory(id); err != nil {
		rsp.Code = common.ResponseCodeInternalErr
		rsp.Message = "删除产品类目失败"
		err = fmt.Errorf("model.DeleteProductCategory[id=%d]: %w", id, err)
		return
	}

	return
}

// UpdateProductCategory 更新产品类目逻辑
func UpdateProductCategory(id, parentID int64, categoryName, categoryNameEN, image, desc, descEN string) (rsp *common.Response, err error) {
	rsp = common.NewOKResponse()

	cate := &model.ProductCategory{
		ID:             id,
		ParentID:       parentID,
		CategoryName:   categoryName,
		CategoryNameEN: categoryNameEN,
		Image:          image,
		Detail:         desc,
		DetailEN:       descEN,
	}

	if err = model.UpdateProductCategory(cate); err != nil {
		rsp.Code = common.ResponseCodeInternalErr
		rsp.Message = "更新产品类目失败"
		err = fmt.Errorf("model.UpdateProductCategory[category=%+v]: %w", *cate, err)
		return
	}

	return
}

// QueryProductCategoryList 查询产品类目列表逻辑
func QueryProductCategoryList(parentID int64) (rsp *common.Response, err error) {
	rsp = common.NewOKResponse()

	list, err := model.QueryProductCategoryList(parentID)
	if err != nil {
		rsp.Code = common.ResponseCodeInternalErr
		rsp.Message = "查询产品类目列表失败"
		err = fmt.Errorf("model.QueryProductCategoryList[parent_id=%d]: %w", parentID, err)
		return
	}

	rsp.Data = common.ResponseData{
		"list": list,
	}

	return
}
