package main

import (
	"context"
	"log"

	"github.com/ItsXomyak/scam-list/config"
	"github.com/ItsXomyak/scam-list/internal/adapter/postgres"
	"github.com/ItsXomyak/scam-list/internal/domain/entity"
	pgclient "github.com/ItsXomyak/scam-list/pkg/postgres"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New(".")
	if err != nil {
		log.Fatal(err)
	}

	client, err := pgclient.New(ctx, cfg.Postgres.GetDsn(), &pgclient.Postgres{})
	if err != nil {
		log.Fatal(err)
	}

	domRepo := postgres.NewDomain(client.Pool)

	// тестовые значения для всех полей
	company := "Example Ltd."
	country := "US"
	scamType := "phishing"
	verifiedBy := "Admin"
	verificationMethod := "manual"
	riskScore := "85.5"

	res, err := domRepo.CreateDomain(ctx, entity.CreateDomainParams{
		Domain:             "example.com",
		Status:             "scam", // "verified", "scam", "suspicious"
		CompanyName:        &company,
		Country:            &country,
		ScamSources:        []string{"source1", "source2"},
		ScamType:           &scamType,
		VerifiedBy:         &verifiedBy,
		VerificationMethod: &verificationMethod,
		RiskScore:          &riskScore,
		Reasons:            []string{"suspicious activity", "reported by users"},
		Metadata: [][]byte{
			[]byte(`{"module":"dns-check","result":"blacklisted"}`),
			[]byte(`{"module":"whois","result":"hidden registrar"}`),
		},
	})
	if err != nil {
		log.Fatal("failed to create domain:", err)
	}

	log.Printf("created domain: %+v\n", res)
}
