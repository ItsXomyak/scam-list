package dto

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/ItsXomyak/scam-list/internal/domain/entity"
	"github.com/ItsXomyak/scam-list/pkg/validator"
)

var (
	ValidDomainStatus = []string{"verified", "scam", "suspicious"}
)

func ValidateCreateDomain(v *validator.Validator, r *entity.CreateDomainParams) {
	validateDomain(v, r.Domain)      // required
	validateStatusValue(v, r.Status) // required

	validateCompanyName(v, r.CompanyName)               // optional
	validateCountry(v, r.Country)                       // optional
	validateScamSources(v, r.ScamSources)               // optional
	validateScamType(v, r.ScamType)                     // optional
	validateVerifiedBy(v, r.VerifiedBy)                 // optional
	validateVerificationMethod(v, r.VerificationMethod) // optional
	validateRiskScore(v, r.RiskScore)                   // optional
	validateReasons(v, r.Reasons)                       // optional
	validateMetadata(v, r.Metadata)                     // optional
}

func ValidatePatchDomain(v *validator.Validator, r *entity.Domain) {
	validateStatusValue(v, r.Status) // required

	validateCompanyName(v, r.CompanyName)               // optional
	validateCountry(v, r.Country)                       // optional
	validateScamSources(v, r.ScamSources)               // optional
	validateScamType(v, r.ScamType)                     // optional
	validateVerifiedBy(v, r.VerifiedBy)                 // optional
	validateVerificationMethod(v, r.VerificationMethod) // optional
	validateRiskScore(v, r.RiskScore)                   // optional
	validateReasons(v, r.Reasons)                       // optional
	validateMetadata(v, r.Metadata)                     // optional
}

/* -------------------- Field helpers -------------------- */

func validateDomain(v *validator.Validator, domain string) {
	v.Check(domain != "", "domain", "must be provided")
	v.Check(len(domain) <= 253, "domain", "must be at most 253 characters")
	v.Check(IsValidDomainName(domain), "domain", "must be a valid domain name (e.g., example.com)")
}

func validateStatusValue(v *validator.Validator, status string) {
	v.Check(status != "", "status", "must be provided")
	v.Check(validator.PermittedValue(status, ValidDomainStatus...), "status",
		fmt.Sprintf("invalid status, available: %s", strings.Join(ValidDomainStatus, ", ")))
}

func validateCompanyName(v *validator.Validator, name *string) {
	if name == nil {
		return
	}
	v.Check(len(*name) != 0, "company_name", "must be provided")
	v.Check(len(*name) <= 255, "company_name", "must be at most 255 characters")
}

func validateCountry(v *validator.Validator, country *string) {
	if country == nil {
		return
	}
	v.Check(len(*country) == 2, "country", "must be exactly 2 characters (ISO 3166-1 alpha-2)")
}

func validateScamSources(v *validator.Validator, sources []string) {
	if sources == nil {
		return
	}
	v.Check(len(sources) != 0, "scam_sources", "must be provided")
	for i, s := range sources {
		field := fmt.Sprintf("scam_sources[%d]", i)
		v.Check(s != "", field, "must be provided")
		v.Check(len(s) <= 100, field, "must be at most 100 characters")
	}
}

func validateScamType(v *validator.Validator, t *string) {
	if t == nil {
		return
	}
	v.Check(len(*t) != 0, "scam_type", "must be provided")
	v.Check(len(*t) <= 100, "scam_type", "must be at most 100 characters")
}

func validateVerifiedBy(v *validator.Validator, by *string) {
	if by == nil {
		return
	}
	v.Check(len(*by) != 0, "verified_by", "must be provided")
	v.Check(len(*by) <= 100, "verified_by", "must be at most 100 characters")
}

func validateVerificationMethod(v *validator.Validator, m *string) {
	if m == nil {
		return
	}
	v.Check(len(*m) != 0, "verification_method", "must be provided")
	v.Check(len(*m) <= 100, "verification_method", "must be at most 100 characters")
}

func validateRiskScore(v *validator.Validator, score *float64) {
	if score == nil {
		return
	}
	val := *score
	v.Check(val >= 0 && val <= 100, "risk_score", "must be between 0 and 100 inclusive")
	v.Check(decimalPlaces(val) <= 2, "risk_score", "must have at most 2 decimal places")
}

func validateReasons(v *validator.Validator, reasons []string) {
	if reasons == nil {
		return
	}
	v.Check(len(reasons) != 0, "reasons", "must be provided")
	for i, s := range reasons {
		field := fmt.Sprintf("reasons[%d]", i)
		v.Check(strings.TrimSpace(s) != "", field, "must not be empty")
	}
}

func validateMetadata(v *validator.Validator, meta []json.RawMessage) {
	if meta == nil {
		return
	}
	v.Check(len(meta) != 0, "metadata", "must be provided")
	for i, raw := range meta {
		var tmp any
		if err := json.Unmarshal(raw, &tmp); err != nil {
			field := fmt.Sprintf("metadata[%d]", i)
			v.Check(false, field, "must be valid JSON")
		}
	}
}

// IsValidDomainName validates a DNS hostname (RFC 1035-ish).
// - total length <= 253 (already checked)
// - labels 1..63 chars, letters/digits/hyphen, cannot start/end with hyphen
// - at least one dot
func IsValidDomainName(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasSuffix(s, ".") { // disallow trailing dot (FQDN) for simplicity
		return false
	}
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}
	labelRe := regexp.MustCompile(`(?i)^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)
	for _, p := range parts {
		if len(p) < 1 || len(p) > 63 {
			return false
		}
		if !labelRe.MatchString(p) {
			return false
		}
	}
	return true
}

func decimalPlaces(f float64) int {
	// Scale and compare to rounded to avoid FP noise.
	for dp := 0; dp <= 6; dp++ { // cap to 6 to avoid loops
		scale := math.Pow(10, float64(dp))
		if math.Abs(math.Round(f*scale)-(f*scale)) < 1e-9 {
			return dp
		}
	}
	return 7 // fallback => more than 6
}
