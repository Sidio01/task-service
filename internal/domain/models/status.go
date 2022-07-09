package models

type StatusApproved struct {
	Status string `json:"status" example:"approved"`
}

type StatusDeclined struct {
	Status string `json:"status" example:"declined"`
}

type StatusDeleted struct {
	Status string `json:"status" example:"deleted"`
}
