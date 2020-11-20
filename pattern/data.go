package pattern

type Data struct {
	Helpers
	Data string
}

// JQ 执行 jq 表达式查询 json 字符串，如果表达式错误，则返回空字符串
func (d Data) JQ(expression string, data ...string) string {
	if len(data) > 0 {
		return d.jq(data[0], expression, true)
	}

	return d.jq(d.Data, expression, true)
}

// JQE 执行 jq 表达式查询 json 字符串，如果表达式错误，则返回 `<ERORR> 错误详情`
func (d Data) JQE(expression string, data ...string) string {
	if len(data) > 0 {
		return d.jq(data[0], expression, false)
	}

	return d.jq(d.Data, expression, false)
}

// DOMOne 从 HTML DOM 对象中查询第 index 个匹配 selector 的元素内容
func (d Data) DOMOne(selector string, index int, data ...string) string {
	if len(data) > 0 {
		return d.domOne(selector, index, data[0])
	}

	return d.domOne(selector, index, d.Data)
}

// DOM 从 HTML DOM 对象中查询所有匹配 selector 的元素
func (d Data) DOM(selector string, data ...string) []string {
	if len(data) > 0 {
		return d.dom(selector, data[0])
	}

	return d.dom(selector, d.Data)
}
