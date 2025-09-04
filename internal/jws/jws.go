package jws

type JWSValidator interface {
	Validate(header string, payload string, signature string) error
}
