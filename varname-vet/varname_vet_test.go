package main

import "slices"
import "testing"
import "golang.org/x/tools/go/analysis/analysistest"

func TestVarnameVet(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer, "clean", "long", "generated")
}

func TestWords(t *testing.T) {
	var tests = []struct {
		name string
		want []string
	}{
		{"a", []string{"a"}},
		{"config", []string{"config"}},
		{"Config", []string{"Config"}},
		{"activeAddr", []string{"active", "Addr"}},
		{"userConfigParam", []string{"user", "Config", "Param"}},
		{"parsedJSONDocument", []string{"parsed", "JSON", "Document"}},
		{"ID", []string{"ID"}},
		{"userID", []string{"user", "ID"}},
		{"userIDs", []string{"user", "IDs"}},
		{"IDsByName", []string{"IDs", "By", "Name"}},
		{"HTTPServer", []string{"HTTP", "Server"}},
		{"HTTP2Server", []string{"HTTP2", "Server"}},
		{"isASet", []string{"is", "A", "Set"}},
		{"utf8Reader", []string{"utf8", "Reader"}},
		{"sha256", []string{"sha256"}},
		{"base64URLEncoding", []string{"base64", "URL", "Encoding"}},
		{"x2y", []string{"x2y"}},
		{"user_config_param", []string{"user", "config", "param"}},
		{"MAX_RETRIES", []string{"MAX", "RETRIES"}},
		{"IPv4Addr", []string{"I", "Pv4", "Addr"}},
		{"_", nil},
		{"_tmp", []string{"tmp"}},
		{"a__b", []string{"a", "b"}},
		{"größeWert", []string{"größe", "Wert"}},
	}
	for _, tt := range tests {
		if got := words(tt.name); !slices.Equal(got, tt.want) {
			t.Errorf("words(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
