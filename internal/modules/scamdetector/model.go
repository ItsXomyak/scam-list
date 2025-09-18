package scamdetector

type ScamDetectorRes struct {
	TotalPercent    float64
	DomainAge       string
	DomainDate      string
	Blacklist       string
	Https           bool
	ProximityToScam bool
	Description     string //div about-text
	Country         string
	CountryCode     string
	Email           string
	SslDate         string
	SslIssuer       string
	ServerName      string
	Review          string
}
