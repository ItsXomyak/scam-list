package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DomainDB struct {
	Domain             string
	Status             string
	CompanyName        sql.NullString
	Country            sql.NullString
	ScamSources        pq.StringArray
	ScamType           sql.NullString
	VerifiedAt         sql.NullTime
	VerifiedBy         sql.NullString
	VerificationMethod sql.NullString
	ExpiresAt          sql.NullTime
	RiskScore          sql.NullString
	Reasons            pq.StringArray
	Metadata           pq.ByteaArray
	CreatedAt          sql.NullTime
	UpdatedAt          sql.NullTime
	LastCheckAt        sql.NullTime
}

type PendingModerationDB struct {
	Domain         string
	CheckID        uuid.UUID
	Reasons        pq.StringArray
	SourceModules  pq.StringArray
	Priority       sql.NullInt32
	Status         sql.NullString
	AssignedTo     sql.NullString
	SubmittedAt    sql.NullTime
	ResolvedAt     sql.NullTime
	ModeratorNotes  sql.NullString
	CreatedAt      sql.NullTime
}