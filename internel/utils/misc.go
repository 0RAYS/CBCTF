package utils

// TidyPaginate 实现内存分页逻辑, 处理各类边界场景
func TidyPaginate(length, limit, offset int) (int, int) {
	if length == 0 {
		return 0, 0
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= length {
		return 0, 0
	}
	var end int
	if limit < 0 {
		end = length // limit<=0 时视为取全部剩余数据
	} else {
		end = offset + limit
	}
	if end > length {
		end = length
	}
	return offset, end
}

func Ptr[T any](v T) *T {
	return &v
}
