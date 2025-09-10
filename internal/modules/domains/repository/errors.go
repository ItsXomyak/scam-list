package repository

import "errors"

var (
	ErrDomainAlreadyExists = errors.New("domain already exists")
	ErrDomainNotFound      = errors.New("domain not found")
	ErrDomainNotVerified   = errors.New("domain not verified")
	ErrDomainNotCreated    = errors.New("domain not created")
)
