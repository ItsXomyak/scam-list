package server

import "github.com/ItsXomyak/scam-list/internal/transport/http/handler"

type Verifier interface {
	handler.Verifier
}
