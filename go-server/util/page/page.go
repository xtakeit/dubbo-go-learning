package page

const DefaultPageSize = 20

// 转化分页参数为数据库查询Limit
func PageToLimit(pageNum, pageSize int64) (from, size int64) {
	size = pageSize
	if size <= 0 {
		size = DefaultPageSize
	}
	
	from = (pageNum - 1) * size
	if from <= 0 {
		from = 0
	}

	return
}
