package models

type RunTask struct {
	ApprovalLogins []string `json:"approvalLogins" swaggertype:"array,string" example:"test626,zxcvb"`
	InitiatorLogin string   `json:"initiatorLogin" example:"test123"`
	Name           string   `json:"name" example:"test task"`
	Text           string   `json:"text" example:"this is test task 1"`
}
