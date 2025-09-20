package server

import "github.com/ItsXomyak/scam-list/internal/adapter/http/handler"

type Verifier interface {
	handler.Verifier
}

type DomainService interface {
	handler.DomainRepository
}
