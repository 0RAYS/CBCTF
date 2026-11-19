package utils

// TidyPaginate 把 offset/limit 规范成内存分页的起止位置，limit<=0 表示取全部剩余数据
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
