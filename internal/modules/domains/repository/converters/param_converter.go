package converters

import (
	"github.com/sqlc-dev/pqtype"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
	"github.com/ItsXomyak/scam-list/internal/modules/domains/repository"
)

func ToDBCreateDomainParams(params *entity.CreateDomainParams) *repository.СreateDomainParams {
	return &repository.СreateDomainParams{
		Domain:             params.Domain,
		CompanyName:        stringToSqlNullString(params.CompanyName),
		Country:            stringToSqlNullString(params.Country),
		ScamSources:        params.ScamSources,
		ScamType:           stringToSqlNullString(params.ScamType),
		VerifiedAt:         timeToSqlNullTime(params.VerifiedAt),
		VerifiedBy:         stringToSqlNullString(params.VerifiedBy),
		VerificationMethod: stringToSqlNullString(params.VerificationMethod),
		ExpiresAt:          timeToSqlNullTime(params.ExpiresAt),
		RiskScore:          stringToSqlNullString(params.RiskScore),
		Reasons:            params.Reasons,
		Metadata:           byteSliceToNullRawMessage(params.Metadata),
		LastCheckAt:        timeToSqlNullTime(params.LastCheckAt),
	}
}

// GetDomainsByRiskScoreParams
func ToDBGetDomainsByRiskScoreParams(params *entity.GetDomainsByRiskScoreParams) repository.GetDomainsByRiskScoreParams {
	return repository.GetDomainsByRiskScoreParams{
		RiskScore:   stringToSqlNullString(params.RiskScore),
		RiskScore_2: stringToSqlNullString(params.RiscScore2), // Обратите внимание на разницу в именах
		Limit:       params.Limit,
		Offset:      params.Offset,
	}
}

// GetDomainsByStatusParams (добавьте в entity если нужно)
func ToDBGetDomainsByStatusParams(status string, limit, offset int32) repository.GetDomainsByStatusParams {
	return repository.GetDomainsByStatusParams{
		Status: status,
		Limit:  limit,
		Offset: offset,
	}
}

// GetDomainsForRecheckParams
func ToDBGetDomainsForRecheckParams(params *entity.GetDomainsForRecheckParams) repository.GetDomainsForRecheckParams {
	return repository.GetDomainsForRecheckParams{
		LastCheckAt: timeToSqlNullTime(params.LastCheckAt),
		Limit:       params.Limit,
	}
}

// MarkDomainAsScamParams
func ToDBMarkDomainAsScamParams(params *entity.MarkDomainAsScamParams) repository.MarkDomainAsScamParams {
	return repository.MarkDomainAsScamParams{
		Domain:      params.Domain,
		ScamSources: params.ScamSources,
		ScamType:    stringToSqlNullString(params.ScamType),
		RiskScore:   stringToSqlNullString(params.RiskScore),
		Reasons:     params.Reasons,
	}
}

func ToDBUpdateDomainStatusParams(params *entity.UpdateDomainStatusParams) repository.UpdateDomainStatusParams {
	return repository.UpdateDomainStatusParams{
		Domain:    params.Domain,
		Status:    params.Status,
		RiskScore: stringToSqlNullString(params.RiskScore),
		Reasons:   params.Reasons,
	}
}

func ToDBVerifyDomainParams(params *entity.VerifyDomainParams) repository.VerifyDomainParams {
	return repository.VerifyDomainParams{
		Domain:             params.Domain,
		VerifiedAt:         timeToSqlNullTime(params.VerifiedAt),
		VerifiedBy:         stringToSqlNullString(params.VerifiedBy),
		VerificationMethod: stringToSqlNullString(params.VerificationMethod),
		ExpiresAt:          timeToSqlNullTime(params.ExpiresAt),
		RiskScore:          stringToSqlNullString(params.RiskScore),
		Reasons:            params.Reasons,
	}
}

// Вспомогательные функции
func byteSliceToNullRawMessage(data [][]byte) pqtype.NullRawMessage {
	if data != nil && len(data) > 0 {
		// TODO: Здесь нужно определить логику преобразования [][]byte в RawMessage
		// В зависимости от того, как хранитcя metadata
		if len(data) == 1 {
			return pqtype.NullRawMessage{
				RawMessage: data[0],
				Valid:      true,
			}
		}
		// Или сериализовать массив в JSON
	}
	return pqtype.NullRawMessage{Valid: false}
}

// // Обратные конвертеры (если понадобятся)
// func ToEntityCreateDomainParams(params *repository.CreateDomainParams) *entity.CreateDomainParams {
// 	return &entity.CreateDomainParams{
// 		Domain:             params.Domain,
// 		CompanyName:        sqlNullStringToPtr(params.CompanyName),
// 		Country:            sqlNullStringToPtr(params.Country),
// 		ScamSources:        params.ScamSources,
// 		ScamType:           sqlNullStringToPtr(params.ScamType),
// 		VerifiedAt:         sqlNullTimeToPtr(params.VerifiedAt),
// 		VerifiedBy:         sqlNullStringToPtr(params.VerifiedBy),
// 		VerificationMethod: sqlNullStringToPtr(params.VerificationMethod),
// 		ExpiresAt:          sqlNullTimeToPtr(params.ExpiresAt),
// 		RiskScore:          sqlNullStringToPtr(params.RiskScore),
// 		Reasons:            params.Reasons,
// 		Metadata:           nullRawMessageToByteSlice(params.Metadata),
// 		LastCheckAt:        sqlNullTimeToPtr(params.LastCheckAt),
// 	}
// }

func nullRawMessageToByteSlice(msg pqtype.NullRawMessage) [][]byte {
	if msg.Valid {
		return [][]byte{msg.RawMessage}
	}
	return nil
}

