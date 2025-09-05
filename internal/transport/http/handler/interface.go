package handler

type Verifier interface {
	VerifyDomain(domain string) error
}
