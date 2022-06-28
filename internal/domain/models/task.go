package models

import "database/sql"

type RunTask struct { // TODO: добавить теги для базы данных
	ApprovalLogins []string `json:"approvalLogins"`
	InitiatorLogin string   `json:"initiatorLogin"`
}

type Approval struct {
	// Approved      bool   `json:"approved"`
	// Sent          bool   `json:"sent"`
	Approved      sql.NullBool `json:"approved"`
	Sent          sql.NullBool `json:"sent"`
	N             int          `json:"n"`
	ApprovalLogin string       `json:"approvalLogin"`
}

// func (a *Approval) ChangeApprovedStatus(b bool) {
// 	a.Approved = b
// }

func (a *Approval) ChangeApprovedStatus(b sql.NullBool) {
	a.Approved = b
}

type Task struct {
	UUID           string      `json:"uuid"`
	Name           string      `json:"name"`
	Text           string      `json:"text"`
	InitiatorLogin string      `json:"initiatorLogin"`
	Status         string      `json:"status"`
	Approvals      []*Approval `json:"approvals"`
}
