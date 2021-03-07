package repo

import "errors"

// ErrNotFound 查询元素不存在
var ErrNotFound = errors.New("no such element")

// 标签，用于规则匹配
type Labels map[string]string
