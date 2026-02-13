package loglint

import (
	"fmt"
	"regexp"
	"strings"
)

type OptionsBuilder struct {
	opts *Options
	err  error
}

func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{opts: DefaultOptions()}
}

func (b *OptionsBuilder) SetStartsWithLowercase(enabled bool) *OptionsBuilder {
	b.opts.StartsWithLowercase = enabled
	return b
}

func (b *OptionsBuilder) SetEngLettersOnly(enabled bool) *OptionsBuilder {
	b.opts.EngLettersOnly = enabled
	return b
}

func (b *OptionsBuilder) SetNoSpecialSymbols(enabled bool) *OptionsBuilder {
	b.opts.NoSpecialSymbols = enabled
	return b
}

func (b *OptionsBuilder) SetAllowedSymbols(symbols string) *OptionsBuilder {
	b.opts.AllowedSymbols = []rune(symbols)

	return b
}

func (b *OptionsBuilder) SetSensitiveKeywords(keywords []string) *OptionsBuilder {
	if b.err != nil {
		return b
	}

	b.opts.SensitiveKeywordsRe = nil

	for _, kw := range keywords {
		if strings.TrimSpace(kw) == "" {
			b.err = fmt.Errorf("sensitive keyword must not be empty")
			return b
		}

		re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(strings.TrimSpace(kw)))
		if err != nil {
			b.err = err
			return b
		}

		b.opts.SensitiveKeywordsRe = append(b.opts.SensitiveKeywordsRe, re)
	}

	return b
}

func (b *OptionsBuilder) SetSensitivePatterns(patterns []string) *OptionsBuilder {
	if b.err != nil {
		return b
	}

	b.opts.SensitivePatterns = nil

	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			b.err = err
			return b
		}

		b.opts.SensitivePatterns = append(b.opts.SensitivePatterns, re)
	}

	return b
}

func (b *OptionsBuilder) Build() (*Options, error) {
	if b.err != nil {
		return nil, b.err
	}

	return cloneOptions(b.opts), nil
}

func cloneOptions(src *Options) *Options {
	if src == nil {
		return nil
	}

	dst := *src
	dst.AllowedSymbols = append([]rune(nil), src.AllowedSymbols...)
	dst.SensitiveKeywordsRe = append([]*regexp.Regexp(nil), src.SensitiveKeywordsRe...)
	dst.SensitivePatterns = append([]*regexp.Regexp(nil), src.SensitivePatterns...)

	return &dst
}
