package utils

import "testing"

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{
			in:   "https://chatgpt.com/g/g-p-68c01bfbb8ac8191837bfe679e8102bf-temutjin",
			want: "chatgpt.com",
		},
		{
			in:   "chatgpt.com",
			want: "chatgpt.com",
		},
		{
			in:   "  chatgpt.com  ",
			want: "chatgpt.com",
		},
		{
			in:   "HTTP://EXAMPLE.COM:8080/path?q=1",
			want: "example.com",
		},
		{
			in:   "//example.com/path",
			want: "example.com",
		},
		{
			in:   "https://sub.domain.example.",
			want: "sub.domain.example",
		},
		{
			in:   "http://127.0.0.1:3000",
			want: "127.0.0.1",
		},
		{
			in:   "http://[::1]:8080",
			want: "::1",
		},
		{
			in:      "",
			wantErr: true,
		},
		{
			in:      "://",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		got, err := ExtractDomain(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ExtractDomain(%q) expected error, got none", tc.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("ExtractDomain(%q) unexpected error: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ExtractDomain(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
}
