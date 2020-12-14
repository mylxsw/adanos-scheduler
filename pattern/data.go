package pattern

import "github.com/tidwall/gjson"

// Data 封装了助手函数和匹配数据的结构体
type Data struct {
	Helpers
	Data string
}

// JQ 执行 JQuery 表达式查询 json 字符串，如果表达式错误，则返回空字符串
func (d Data) JQ(expression string, data ...string) string {
	if len(data) > 0 {
		return d.JQuery(data[0], expression, true)
	}

	return d.JQuery(d.Data, expression, true)
}

// JQE 执行 JQuery 表达式查询 json 字符串，如果表达式错误，则返回 `<ERORR> 错误详情`
func (d Data) JQE(expression string, data ...string) string {
	if len(data) > 0 {
		return d.JQuery(data[0], expression, false)
	}

	return d.JQuery(d.Data, expression, false)
}

// DOMOne 从 HTML DOM 对象中查询第 index 个匹配 selector 的元素内容
func (d Data) DOMOne(selector string, index int, data ...string) string {
	if len(data) > 0 {
		return d.DOMQueryOne(selector, index, data[0])
	}

	return d.DOMQueryOne(selector, index, d.Data)
}

// DOM 从 HTML DOM 对象中查询所有匹配 selector 的元素
func (d Data) DOM(selector string, data ...string) []string {
	if len(data) > 0 {
		return d.DOMQuery(selector, data[0])
	}

	return d.DOMQuery(selector, d.Data)
}

// CtxJSONArray return array from json
func (d Data) CtxJSONArray(path string) []gjson.Result {
	return d.Helpers.JSONArray(d.Data, path)
}

// CtxJSONStrArray get string array from json
func (d Data) CtxJSONStrArray(path string) []string {
	return d.Helpers.JSONStrArray(d.Data, path)
}

// CtxJSONIntArray return int array from json
func (d Data) CtxJSONIntArray(path string) []int64 {
	return d.Helpers.JSONIntArray(d.Data, path)
}

// CtxJSONFloatArray return float64 array from json
func (d Data) CtxJSONFloatArray(path string) []float64 {
	return d.Helpers.JSONFloatArray(d.Data, path)
}

// CtxJSONBoolArray return bool array from json
func (d Data) CtxJSONBoolArray(path string) []bool {
	return d.Helpers.JSONBoolArray(d.Data, path)
}

// CtxJSON return string from json
func (d Data) CtxJSON(path string) string {
	return d.Helpers.JSON(d.Data, path)
}

// CtxJSONInt return int from json
func (d Data) CtxJSONInt(path string) int64 {
	return d.Helpers.JSONInt(d.Data, path)
}

// CtxJSONFloat return float64 from json
func (d Data) CtxJSONFloat(path string) float64 {
	return d.Helpers.JSONFloat(d.Data, path)
}

// CtxJSONBool return bool from json
func (d Data) CtxJSONBool(path string) bool {
	return d.Helpers.JSONBool(d.Data, path)
}
