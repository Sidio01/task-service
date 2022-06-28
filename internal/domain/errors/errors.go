package errors

import "errors"

var (
	ErrInvalidJsonBody          = errors.New("{\"error\": \"invalid json body\"}")
	ErrIdNotFound               = errors.New("{\"error\": \"id not found\"}")
	ErrLoginNotFoundInApprovals = errors.New("{\"error\": \"login not found in approvals\"}")
	ErrMock                     = errors.New("{\"status\": \"error\"}")
)

// TODO: обертка для ошибок пакетов
