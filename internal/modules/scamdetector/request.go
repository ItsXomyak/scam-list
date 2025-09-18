package scamdetector

import (
	"context"
	"net/http"
)

func Request(ctx context.Context, domain string) {
	http.Get("https://www.scam-detector.com/validator/youtube-com-review")
}