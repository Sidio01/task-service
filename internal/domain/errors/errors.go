package errors

import "errors"

var (
	ErrInvalidJsonBody                  = errors.New("{\"error\": \"invalid json body\"}")
	ErrIdNotFound                       = errors.New("{\"error\": \"id not found\"}")
	ErrLoginNotFoundInApprovals         = errors.New("{\"error\": \"login not found in approvals\"}")
	ErrAuthFailed                       = errors.New("{\"error\": \"authorization failed, wrong token\"}")
	ErrTokenLoginNotEqualInitiatorLogin = errors.New("{\"error\": \"token login not equal initiator login\"}")
	ErrMock                             = errors.New("{\"status\": \"error\"}")
)

// TODO: обертка для ошибок пакетов
