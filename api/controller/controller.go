package controller

// IDResponse 用于响应返回结果只有一个 id 的请求
type IDResponse struct {
	ID string `json:"id"`
}