package structs

type JobsCreate struct {
	FileID string `json:"file_id" validate:"required,uuid4"` 
	Type string `json:"job_type" validate:"required,uuid4"` 
	Status string `json:"status" validate:"required,uuid4"` 
}

type JobsQualityControlRequest struct {
	FileID  string `json:"file_id" validate:"required"`
	JobName string `json:"job_name" validate:"omitempty"`
	Type string `json:"job_type" validate:"required"`
}