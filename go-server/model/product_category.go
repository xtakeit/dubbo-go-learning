// Author: Steve Zhang
// Date: 2020/10/18 1:17 下午

package model

import (
	"fmt"

	"go-server/component"
)

// ProductCategory 定义产品类目结构
type ProductCategory struct {
	ID             int64  `db:"id" json:"id"`
	ParentID       int64  `db:"parent_id" json:"parent_id"`
	CategoryName   string `db:"category_name" json:"category_name"`
	CategoryNameEN string `db:"category_name_en" json:"category_name_en"`
	Image          string `db:"image" json:"image"`
	Detail         string `db:"detail" json:"detail"`
	DetailEN       string `db:"detail_en" json:"detail_en"`
	IsDeleted      string `db:"is_deleted" json:"is_deleted"`
	CreatedAt      string `db:"created_at" json:"created_at"`
	UpdatedAt      string `db:"updated_at" json:"updated_at"`
}

const addProductCategorySQL = `INSERT INTO t_product_category (parent_id, category_name, 
category_name_en, image, detail, detail_en) VALUES (?, ?, ?, ?, ?, ?)`

// AddProductCategory 新增产品类目
func AddProductCategory(cate *ProductCategory) (id int64, err error) {
	if _, id, err = component.DBContainer.Exec(
		addProductCategorySQL, cate.ParentID, cate.CategoryName, cate.CategoryNameEN,
		cate.Image, cate.Detail, cate.DetailEN,
	); err != nil {
		err = fmt.Errorf("component.DBContainer.Exec[sql=%s]: %w", addProductCategorySQL, err)
		return
	}

	return
}

const deleteProductCategorySQL = `UPDATE t_product_category SET is_deleted = 1 WHERE id = ?`

// DeleteProductCategory 删除产品类目
func DeleteProductCategory(id int64) (err error) {
	if _, _, err = component.DBContainer.Exec(deleteProductCategorySQL, id); err != nil {
		err = fmt.Errorf("component.DBContainer.Exec[sql=%s]: %w", deleteProductCategorySQL, err)
		return
	}

	return
}

const updateProductCategorySQL = `UPDATE t_product_category SET parent_id = ?, category_name = ?, category_name_en = ?, 
image = ?, detail = ?, detail_en = ? WHERE id = ?`

// UpdateProductCategory 更新产品类目
func UpdateProductCategory(cate *ProductCategory) (err error) {
	if _, _, err = component.DBContainer.Exec(
		updateProductCategorySQL, cate.ParentID, cate.CategoryName, cate.CategoryNameEN, cate.Image,
		cate.Detail, cate.DetailEN, cate.ID,
	); err != nil {
		err = fmt.Errorf("component.DBContainer.Exec[sql=%s]: %w", updateProductCategorySQL, err)
		return
	}

	return
}

const queryProductCategoryListSQL = `SELECT id, parent_id, category_name, category_name_en, image, detail,
detail_en, is_deleted, created_at, updated_at FROM t_product_category WHERE is_deleted = 0 AND parent_id = ?`

// QueryProductCategoryList 查询产品类目列表
func QueryProductCategoryList(parentID int64) (list []*ProductCategory, err error) {
	if err = component.DBContainer.Query(queryProductCategoryListSQL, &list, parentID); err != nil {
		err = fmt.Errorf("component.DBContainer.Query[sql=%s]: %w", queryProductCategoryListSQL, err)
		return
	}
	return
}
