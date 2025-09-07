package converters

import (
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
	"github.com/ItsXomyak/scam-list/internal/modules/domains/repository/models"
)

func ToDomainEntity(dbDomain *models.DomainDB) *entity.Domain {
	if dbDomain == nil {
		return nil
	}

	return &entity.Domain{
		Domain:             dbDomain.Domain,
		Status:             dbDomain.Status,
		CompanyName:        sqlNullStringToPtr(dbDomain.CompanyName),
		Country:            sqlNullStringToPtr(dbDomain.Country),
		ScamSources:        pqStringArrayToSlice(dbDomain.ScamSources),
		ScamType:           sqlNullStringToPtr(dbDomain.ScamType),
		VerifiedAt:         sqlNullTimeToPtr(dbDomain.VerifiedAt),
		VerifiedBy:         sqlNullStringToPtr(dbDomain.VerifiedBy),
		VerificationMethod: sqlNullStringToPtr(dbDomain.VerificationMethod),
		ExpiresAt:          sqlNullTimeToPtr(dbDomain.ExpiresAt),
		RiskScore:          sqlNullStringToPtr(dbDomain.RiskScore),
		Reasons:            pqStringArrayToSlice(dbDomain.Reasons),
		Metadata:           pqByteaArrayToSlice(dbDomain.Metadata),
		CreatedAt:          sqlNullTimeToPtr(dbDomain.CreatedAt),
		UpdatedAt:          sqlNullTimeToPtr(dbDomain.UpdatedAt),
		LastCheckAt:        sqlNullTimeToPtr(dbDomain.LastCheckAt),
	}
}

func ToPendingModerationEntity(dbMod *models.PendingModerationDB) *entity.PendingModeration {
	if dbMod == nil {
		return nil
	}

	return &entity.PendingModeration{
		Domain:         dbMod.Domain,
		CheckID:        dbMod.CheckID.String(),
		Reasons:        pqStringArrayToSlice(dbMod.Reasons),
		SourceModules:  pqStringArrayToSlice(dbMod.SourceModules),
		Priority:       sqlNullInt32ToPtr(dbMod.Priority),
		Status:         sqlNullStringToPtr(dbMod.Status),
		AssignedTo:     sqlNullStringToPtr(dbMod.AssignedTo),
		SubmittedAt:    sqlNullTimeToPtr(dbMod.SubmittedAt),
		ResolvedAt:     sqlNullTimeToPtr(dbMod.ResolvedAt),
		ModeratorNotes: sqlNullStringToPtr(dbMod.ModeratorNotes),
		CreatedAt:      sqlNullTimeToPtr(dbMod.CreatedAt),
	}
}

// вспомогалки для конвертации
func sqlNullStringToPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func sqlNullTimeToPtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func sqlNullInt32ToPtr(i sql.NullInt32) *int32 {
	if i.Valid {
		return &i.Int32
	}
	return nil
}

func pqStringArrayToSlice(arr pq.StringArray) []string {
	if arr == nil {
		return []string{}
	}
	return []string(arr)
}

func pqByteaArrayToSlice(arr pq.ByteaArray) [][]byte {
	if arr == nil {
		return [][]byte{}
	}
	return [][]byte(arr)
}

// Функции для обратной конвертации (если нужно)
func ToDomainDB(domain *entity.Domain) *models.DomainDB {
	if domain == nil {
		return nil
	}

	return &models.DomainDB{
		Domain:             domain.Domain,
		Status:             domain.Status,
		CompanyName:        stringToSqlNullString(domain.CompanyName),
		Country:            stringToSqlNullString(domain.Country),
		ScamSources:        sliceToPqStringArray(domain.ScamSources),
		ScamType:           stringToSqlNullString(domain.ScamType),
		VerifiedAt:         timeToSqlNullTime(domain.VerifiedAt),
		VerifiedBy:         stringToSqlNullString(domain.VerifiedBy),
		VerificationMethod: stringToSqlNullString(domain.VerificationMethod),
		ExpiresAt:          timeToSqlNullTime(domain.ExpiresAt),
		RiskScore:          stringToSqlNullString(domain.RiskScore),
		Reasons:            sliceToPqStringArray(domain.Reasons),
		Metadata:           sliceToPqByteaArray(domain.Metadata),
		CreatedAt:          timeToSqlNullTime(domain.CreatedAt),
		UpdatedAt:          timeToSqlNullTime(domain.UpdatedAt),
		LastCheckAt:        timeToSqlNullTime(domain.LastCheckAt),
	}
}

func stringToSqlNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func timeToSqlNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

func sliceToPqStringArray(slice []string) pq.StringArray {
	if slice == nil {
		return pq.StringArray{}
	}
	return pq.StringArray(slice)
}

func sliceToPqByteaArray(slice [][]byte) pq.ByteaArray {
	if slice == nil {
		return pq.ByteaArray{}
	}
	return pq.ByteaArray(slice)
}