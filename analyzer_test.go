package loglint

import (
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLoglintAnalyzer(t *testing.T) {
	t.Parallel()
	tests := map[string]*Options{
		"starts_with_lowercase":        {StartsWithLowercase: true},
		"eng_letters_only":             {EngLettersOnly: true},
		"no_special_symbols":           {NoSpecialSymbols: true},
		"no_special_symbols_whitelist": {NoSpecialSymbols: true, AllowedSymbols: []rune{'_', '.', ':', '😀'}},
		"multiple_rules": {
			StartsWithLowercase: true,
			EngLettersOnly:      true,
			NoSpecialSymbols:    true,
			SensitiveKeywordsRe: []*regexp.Regexp{
				regexp.MustCompile("(?i)password"),
				regexp.MustCompile("(?i)token"),
			},
			SensitivePatterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\bpassword\b`),
				regexp.MustCompile(`(?i)\btoken\b`),
				regexp.MustCompile(`\b\d{4}-\d{4}-\d{4}-\d{4}\b`),
			},
		},
		"no_sensitive_keywords": {SensitiveKeywordsRe: []*regexp.Regexp{
			regexp.MustCompile("(?i)password"),
			regexp.MustCompile("(?i)ip_address"),
			regexp.MustCompile("(?i)name_of_the_first_pet"),
		}},
		"no_sensitive_patterns": {SensitivePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\bpassword\b`),
			regexp.MustCompile(`(?i)\bsecret\b`),
			regexp.MustCompile(`(?i)\btoken\b`),
			regexp.MustCompile(`(?i)\bip_address\b`),
			regexp.MustCompile(`\b\d{4}-\d{4}-\d{4}-\d{4}\b`),
		}},
	}

	for dir, opt := range tests {
		an := NewAnalyzer(opt)
		testdata := analysistest.TestData()
		analysistest.RunWithSuggestedFixes(t, testdata, an, dir)
	}
}
