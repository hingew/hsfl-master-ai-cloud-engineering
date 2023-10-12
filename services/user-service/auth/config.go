package auth

type Config interface {
	ReadPrivateKey() (any, error)
}
