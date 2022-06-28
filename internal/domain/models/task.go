package models

type RunTask struct {
	ApprovalLogins []string `json:"approvalLogins"`
	InitiatorLogin string   `json:"initiatorLogin"`
}

type Approval struct {
	Approved      bool   `json:"approved"`
	Sent          bool   `json:"sent"`
	N             int    `json:"n"`
	ApprovalLogin string `json:"approvalLogin"`
}

func (a *Approval) ChangeApprovedStatus(b bool) {
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
