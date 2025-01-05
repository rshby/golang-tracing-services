package dto

type GetProductResponseDTO struct {
	Message string `json:"message"`
	TraceID string `json:"traceID"`
	Data    string `json:"data"`
}
