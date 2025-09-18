package scamdetector

type ScamDetectorRes struct {
	// ===== Technical Analysis =====
	DomainAge string `json:"Domain age"`

	CompanyData string `json:"Company Data"`
	Website     string `json:"Website"`
	SSLValid    string `json:"SSL certificate valid"`
	SSLIssuer   string `json:"SSL issuer"`
	WHOISReg    string `json:"WHOIS registration date"`
	WHOISUpd    string `json:"WHOIS last update date"`
	WHOISRenew  string `json:"WHOIS renew date"`

	Owner            string `json:"Owner"`
	Administrator    string `json:"Administrator"`
	TechnicalContact string `json:"Technical Contact"`

	RegistrarName    string `json:"Name"`
	RegistrarIanaID  string `json:"IANA ID"`
	RegistrarWebsite string `json:"Register website"`
	RegistrarEmail   string `json:"E-mail"`
	RegistrarPhone   string `json:"Phone"`

	ServerName []string `json:"Server Name"`

	// ===== Summary =====
	TotalPercent    string `json:"totalPercent"`
	DomainAgeSum    string `json:"domainAge"`
	DomainDate      string `json:"domainDate"`
	BlackList       string `json:"blackList"`
	HttpsConnection string `json:"httpsConnection"`
	SiteDescription string `json:"siteDescription"`
}
