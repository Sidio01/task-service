package ports

type Grpc interface {
	Validate(refreshCookie, accessCookie string) (bool, string, error)
}
