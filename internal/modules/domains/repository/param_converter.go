package repository

import (
	"github.com/sqlc-dev/pqtype"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
)

func ToDBCreateDomainParams(params *entity.CreateDomainParams) CreateDomainParams {
	return CreateDomainParams{
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
func ToDBGetDomainsByRiskScoreParams(params *entity.GetDomainsByRiskScoreParams) GetDomainsByRiskScoreParams {
	return GetDomainsByRiskScoreParams{
		RiskScore:   stringToSqlNullString(params.RiskScore),
		RiskScore_2: stringToSqlNullString(params.RiscScore2), 
		Limit:       params.Limit,
		Offset:      params.Offset,
	}
}

// GetDomainsByStatusParams 
func ToDBGetDomainsByStatusParams(status string, limit, offset int32) GetDomainsByStatusParams {
	return GetDomainsByStatusParams{
		Status: status,
		Limit:  limit,
		Offset: offset,
	}
}

func ToDBGetDomainsByStatusParams2(params *entity.GetDomainsByStatusParams) GetDomainsByStatusParams {
	return GetDomainsByStatusParams{
		Status: params.Status,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
}

// GetDomainsForRecheckParams
func ToDBGetDomainsForRecheckParams(params *entity.GetDomainsForRecheckParams) GetDomainsForRecheckParams {
	return GetDomainsForRecheckParams{
		LastCheckAt: timeToSqlNullTime(params.LastCheckAt),
		Limit:       params.Limit,
	}
}

// MarkDomainAsScamParams
func ToDBMarkDomainAsScamParams(params *entity.MarkDomainAsScamParams) MarkDomainAsScamParams {
	return MarkDomainAsScamParams{
		Domain:      params.Domain,
		ScamSources: params.ScamSources,
		ScamType:    stringToSqlNullString(params.ScamType),
		RiskScore:   stringToSqlNullString(params.RiskScore),
		Reasons:     params.Reasons,
	}
}

func ToDBUpdateDomainStatusParams(params *entity.UpdateDomainStatusParams) UpdateDomainStatusParams {
	return UpdateDomainStatusParams{
		Domain:    params.Domain,
		Status:    params.Status,
		RiskScore: stringToSqlNullString(params.RiskScore),
		Reasons:   params.Reasons,
	}
}

func ToDBVerifyDomainParams(params *entity.VerifyDomainParams) VerifyDomainParams {
	return VerifyDomainParams{
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
// func ToEntityCreateDomainParams(params *CreateDomainParams) *entity.CreateDomainParams {
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

