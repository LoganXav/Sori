package structs

type JobsCreate struct {
	FileID      string `json:"file_id" validate:"required,uuid4"`
	ReferenceID string `json:"reference_id"`
	Name        string `json:"job_name" validate:"required"`
	Type        string `json:"job_type" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

type JobsQualityControlRequest struct {
	FileID string `json:"file_id" validate:"required"`
	Name   string `json:"job_name" validate:"required"`
	Type   string `json:"job_type" validate:"required,eq=quality_control"`
}

type JobsAlignmentRequest struct {
	FileID      string `json:"file_id" validate:"required"`
	ReferenceID string `json:"reference_id" validate:"required"`
	Name        string `json:"job_name" validate:"required"`
	Type        string `json:"job_type" validate:"required,eq=alignment"`
}
