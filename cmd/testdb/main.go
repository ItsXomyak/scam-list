package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ItsXomyak/scam-list/config"
	"github.com/ItsXomyak/scam-list/internal/adapter/postgres"
	"github.com/ItsXomyak/scam-list/internal/domain/entity"
	pgclient "github.com/ItsXomyak/scam-list/pkg/postgres"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New(".env")
	if err != nil {
		log.Fatal(err)
	}

	client, err := pgclient.New(ctx, cfg.Postgres.GetDsn(), &pgclient.Config{
		MaxPoolSize:  5,
		ConnAttempts: 3,
		ConnTimeout:  2_000_000_000, // 2s
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	domRepo := postgres.NewDomain(client.Pool)

	// тестовые значения для всех полей
	company := "Example Ltd."
	country := "US"
	scamType := "phishing"
	verifiedBy := "Admin"
	verificationMethod := "manual"
	riskScore := "85.5"

	// 1. CreateDomain
	res, err := domRepo.CreateDomain(ctx, entity.CreateDomainParams{
		Domain:             "example.com",
		Status:             "scam",
		CompanyName:        &company,
		Country:            &country,
		ScamSources:        []string{"source1", "source2"},
		ScamType:           &scamType,
		VerifiedBy:         &verifiedBy,
		VerificationMethod: &verificationMethod,
		RiskScore:          &riskScore,
		Reasons:            []string{"suspicious activity", "reported by users"},
		Metadata: []json.RawMessage{
			[]byte(`{"module":"dns-check","result":"blacklisted"}`),
			[]byte(`{"module":"whois","result":"hidden registrar"}`),
		},
	})
	if err != nil {
		log.Fatal("failed to create domain:", err)
	}
	prettyPrint("1. created domain", res)

	// 2. GetDomain
	got, err := domRepo.GetDomain(ctx, "example.com")
	if err != nil {
		log.Fatal("failed to get domain:", err)
	}
	prettyPrint("2. get domain", got)

	// 3. GetAllDomains
	all, err := domRepo.GetAllDomains(ctx)
	if err != nil {
		log.Fatal("failed to get all domains:", err)
	}
	log.Println("3. all domains:")
	for _, d := range all {
		prettyPrint("   -> domain", d)
	}

	// 4. UpdateDomain
	got.Status = "suspicious"
	newCompany := "Updated Company"
	got.CompanyName = &newCompany
	got.Reasons = append(got.Reasons, "manual review")

	updated, err := domRepo.UpdateDomain(ctx, got)
	if err != nil {
		log.Fatal("failed to update domain:", err)
	}
	prettyPrint("4. updated domain", updated)

	// 5. GetDomain again
	got2, err := domRepo.GetDomain(ctx, "example.com")
	if err != nil {
		log.Fatal("failed to get domain after update:", err)
	}
	prettyPrint("5. get after update", got2)

	// // 6. DeleteDomain
	// if err := domRepo.DeleteDomain(ctx, "example.com"); err != nil {
	// 	log.Fatal("failed to delete domain:", err)
	// }
	// log.Println("6. domain deleted successfully")
}

// helper: печатает структуру красиво в JSON
func prettyPrint(label string, v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("%s: (error: %v)", label, err)
		return
	}
	log.Printf("%s:\n%s\n", label, data)
}
