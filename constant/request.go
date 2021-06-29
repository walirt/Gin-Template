package constant

type Request struct {
	Version string      `json:"version" binding:"required"`
	Data    interface{} `json:"data" binding:"required"`
}
