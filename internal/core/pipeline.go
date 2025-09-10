package core

import (
	"context"
	"sync"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
)

type ScamChecker interface {
	Check(ctx context.Context, domain string) (*entity.CheckerResult, error) // core функция чекеров с модулей
	Info() string // полная инфа с чекера
}

type DomainService interface {
	GetDomain(ctx context.Context, domain string) (*entity.Domain, error)
}

type DomainPipeline struct {
	checkers []ScamChecker
	domainSvc DomainService
}

func NewDomainPipeline(checkers []ScamChecker, domainSvc DomainService) *DomainPipeline {
	return &DomainPipeline{
		checkers: checkers,
		domainSvc: domainSvc,
	}
}

func (p *DomainPipeline) ProcessDomain(ctx context.Context, url string) (*entity.Domain, error) {
	domain, err := p.domainSvc.GetDomain(ctx, url)
	if err != nil {
		return nil, err // TODO: обработать
	}

	if domain.Domain != "" {
		return domain, nil
	}

	wg := &sync.WaitGroup{}
	resCh := make(chan *entity.CheckerResult, len(p.checkers))
	errCh := make(chan error, len(p.checkers))

	for _, checker := range p.checkers {
		wg.Add(1)
		 go func(checker ScamChecker){
			defer wg.Done()
					result, err := checker.Check(ctx, url)
					if err != nil {
						errCh <- err
						return
					}

					resCh <- result
		 }(checker)
	
	}

	go func () {
		wg.Wait()
		close(resCh)
		close(errCh)
	}()

	total := 0.0
	var results []*entity.CheckerResult
	for {
		select {
		case result, ok := <-resCh:
			if !ok {
				resCh = nil
			} else {
				total += result.TotalScore
				results = append(results, result)
			}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				// TODO: логировать ошибку или собирать их
				_ = err
			}
		}
		
		if resCh == nil && errCh == nil {
			break
		}
	}


	// TODO: доработать формирование ответа с модулей
	verifyResult := &entity.Domain{
		Domain:             url,
		Status:             "unknown", // ну здесь вообще хуйня
		CompanyName:        nil,
		Country:            nil,
		ScamSources:        nil,
		ScamType:           nil,
		VerifiedBy:         nil,
		VerificationMethod: nil,
		RiskScore:          nil,
		Reasons:            nil,
		Metadata:           nil,
		CreatedAt:          nil,
		UpdatedAt:          nil,
	}

	return verifyResult, nil
}


// func (p *DomainPipeline) collectReasons(results []CheckResult) []string {
// 	var reasons []string
// 	for _, result := range results {
// 		if len(result.Reasons) > 0 {
// 			reasons = append(reasons, result.Reasons...)
// 		}
// 	}
// 	return reasons
// }

// func (p *DomainPipeline) collectSources(results []CheckResult) []string {
// 	sources := make([]string, len(results))
// 	for i, result := range results {
// 		sources[i] = result.Source
// 	}
// 	return sources
// }


func stringPtr(s string) *string {
	return &s
}