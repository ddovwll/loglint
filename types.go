package loglint

import (
	"go/token"
	"regexp"

	"golang.org/x/tools/go/analysis"
)

type LogCall struct {
	Msg          string
	MsgIsNotVar  bool
	MsgIsLiteral bool
	FixComputed  bool
	FixAvailable bool
	CombinedFix  analysis.SuggestedFix

	MsgExprPos token.Pos
	MsgExprEnd token.Pos
}

type Options struct {
	StartsWithLowercase bool
	EngLettersOnly      bool
	NoSpecialSymbols    bool
	AllowedSymbols      []rune
	SensitiveKeywordsRe []*regexp.Regexp
	SensitivePatterns   []*regexp.Regexp
}

func DefaultOptions() *Options {
	return &Options{
		StartsWithLowercase: true,
		EngLettersOnly:      true,
		NoSpecialSymbols:    true,
		AllowedSymbols:      nil,
		SensitiveKeywordsRe: nil,
		SensitivePatterns:   nil,
	}
}

type RuleFunc func(pass *analysis.Pass, opts *Options, call *LogCall)
