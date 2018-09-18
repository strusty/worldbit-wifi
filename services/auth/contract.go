package auth

type Auth interface {
	CreateCode(request SendCodeRequest) (string, error)
}
