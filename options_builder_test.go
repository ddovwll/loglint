package loglint

import "testing"

func TestOptionsBuilderBuild(t *testing.T) {
	builder := NewOptionsBuilder().
		SetStartsWithLowercase(false).
		SetEngLettersOnly(false).
		SetNoSpecialSymbols(false).
		SetAllowedSymbols("_.").
		SetSensitiveKeywords([]string{"password", "ip_address"}).
		SetSensitivePatterns([]string{`(?i)\bsecret\b`})

	opts, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if opts.StartsWithLowercase {
		t.Fatalf("expected StartsWithLowercase to be false")
	}
	if opts.EngLettersOnly {
		t.Fatalf("expected EngLettersOnly to be false")
	}
	if opts.NoSpecialSymbols {
		t.Fatalf("expected NoSpecialSymbols to be false")
	}
	if len(opts.AllowedSymbols) != 2 {
		t.Fatalf("expected 2 allowed symbols, got %d", len(opts.AllowedSymbols))
	}
	if !opts.SensitiveKeywordsRe[0].MatchString("PASSWORD") {
		t.Fatalf("expected keyword regex to be case-insensitive")
	}
	if len(opts.SensitivePatterns) != 1 {
		t.Fatalf("expected 1 sensitive pattern")
	}
	if !opts.SensitivePatterns[0].MatchString("secret") {
		t.Fatalf("expected sensitive pattern to match")
	}

	opts.AllowedSymbols[0] = 'x'
	opts2, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected build error on second build: %v", err)
	}
	if opts2.AllowedSymbols[0] != '_' {
		t.Fatalf("expected cloned allowed symbols to be independent")
	}
}

func TestOptionsBuilderRejectsEmptyKeyword(t *testing.T) {
	_, err := NewOptionsBuilder().
		SetSensitiveKeywords([]string{""}).
		Build()
	if err == nil {
		t.Fatalf("expected error for empty sensitive keyword")
	}
}

func TestOptionsBuilderRejectsInvalidPattern(t *testing.T) {
	_, err := NewOptionsBuilder().
		SetSensitivePatterns([]string{"("}).
		Build()
	if err == nil {
		t.Fatalf("expected error for invalid sensitive pattern")
	}
}
