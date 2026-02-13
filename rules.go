package loglint

import (
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var rules = []RuleFunc{
	startsWithLowercase,
	engLettersOnly,
	noSpecialSymbols,
	noSensitiveKeywords,
	noSensitivePatterns,
}

func startsWithLowercase(pass *analysis.Pass, opts *Options, call *LogCall) {
	if !opts.StartsWithLowercase {
		return
	}
	if !call.MsgIsNotVar {
		return
	}

	firstLetter, _ := utf8.DecodeRuneInString(call.Msg)
	if firstLetter >= 'a' && firstLetter <= 'z' {
		return
	}

	d := analysis.Diagnostic{
		Pos:     call.MsgExprPos,
		End:     call.MsgExprEnd,
		Message: "log message must start with a lowercase English letter",
	}

	if fix, ok := combinedFix(opts, call); ok {
		d.SuggestedFixes = append(d.SuggestedFixes, fix)
	}

	pass.Report(d)
}

func engLettersOnly(pass *analysis.Pass, opts *Options, call *LogCall) {
	if !opts.EngLettersOnly {
		return
	}
	if !call.MsgIsNotVar {
		return
	}

	for _, r := range call.Msg {
		if unicode.IsLetter(r) && ((r < 'a' || r > 'z') && (r < 'A' || r > 'Z')) {
			d := analysis.Diagnostic{
				Pos:     call.MsgExprPos,
				End:     call.MsgExprEnd,
				Message: "log message must contain only English letters",
			}

			if fix, ok := combinedFix(opts, call); ok {
				d.SuggestedFixes = append(d.SuggestedFixes, fix)
			}

			pass.Report(d)
			return
		}
	}
}

func noSpecialSymbols(pass *analysis.Pass, opts *Options, call *LogCall) {
	if !opts.NoSpecialSymbols {
		return
	}
	if !call.MsgIsNotVar {
		return
	}

	for _, r := range call.Msg {
		if !unicode.In(r, unicode.Letter, unicode.Space, unicode.Number) && !slices.Contains(opts.AllowedSymbols, r) {
			d := analysis.Diagnostic{
				Pos:     call.MsgExprPos,
				End:     call.MsgExprEnd,
				Message: "found forbidden symbol",
			}

			if fix, ok := combinedFix(opts, call); ok {
				d.SuggestedFixes = append(d.SuggestedFixes, fix)
			}

			pass.Report(d)
			return
		}
	}
}

func noSensitiveKeywords(pass *analysis.Pass, opts *Options, call *LogCall) {
	if len(opts.SensitiveKeywordsRe) == 0 {
		return
	}
	if !call.MsgIsNotVar {
		return
	}

	for _, re := range opts.SensitiveKeywordsRe {
		if re.MatchString(call.Msg) {
			d := analysis.Diagnostic{
				Pos:     call.MsgExprPos,
				End:     call.MsgExprEnd,
				Message: "found sensitive keyword",
			}

			if fix, ok := combinedFix(opts, call); ok {
				d.SuggestedFixes = append(d.SuggestedFixes, fix)
			}

			pass.Report(d)
			return
		}
	}
}

func noSensitivePatterns(pass *analysis.Pass, opts *Options, call *LogCall) {
	if len(opts.SensitivePatterns) == 0 {
		return
	}
	if !call.MsgIsNotVar {
		return
	}

	msg := call.Msg

	for _, pattern := range opts.SensitivePatterns {
		if pattern.MatchString(msg) {
			d := analysis.Diagnostic{
				Pos:     call.MsgExprPos,
				End:     call.MsgExprEnd,
				Message: "found sensitive pattern",
			}

			if fix, ok := combinedFix(opts, call); ok {
				d.SuggestedFixes = append(d.SuggestedFixes, fix)
			}

			pass.Report(d)
			return
		}
	}
}

func combinedFix(opts *Options, call *LogCall) (analysis.SuggestedFix, bool) {
	if call.FixComputed {
		return call.CombinedFix, call.FixAvailable
	}

	call.FixComputed = true
	if !call.MsgIsLiteral {
		return analysis.SuggestedFix{}, false
	}

	fixed := call.Msg
	changed := false

	if len(opts.SensitivePatterns) != 0 {
		if msg, ok := applyNoSensitivePatternsFixText(opts, fixed); ok {
			fixed = msg
			changed = true
		}
	}

	if len(opts.SensitiveKeywordsRe) != 0 {
		if msg, ok := applyNoSensitiveKeywordsFixText(opts, fixed); ok {
			fixed = msg
			changed = true
		}
	}

	if opts.NoSpecialSymbols {
		if msg, ok := applyNoSpecialSymbolsFixText(opts, fixed); ok {
			fixed = msg
			changed = true
		}
	}

	if opts.EngLettersOnly {
		if msg, ok := applyEngLettersOnlyFixText(fixed); ok {
			fixed = msg
			changed = true
		}
	}

	if opts.StartsWithLowercase {
		if msg, ok := applyStartsWithLowercaseFixText(fixed); ok {
			fixed = msg
			changed = true
		}
	}

	if !changed || fixed == call.Msg || len(fixed) == 0 || !isFixResultValid(opts, fixed) {
		return analysis.SuggestedFix{}, false
	}

	call.FixAvailable = true
	call.CombinedFix = analysis.SuggestedFix{
		Message: "apply all available log message fixes",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     call.MsgExprPos,
				End:     call.MsgExprEnd,
				NewText: []byte(strconv.Quote(fixed)),
			},
		},
	}

	return call.CombinedFix, true
}

func isFixResultValid(opts *Options, msg string) bool {
	if opts.StartsWithLowercase {
		r, _ := utf8.DecodeRuneInString(msg)
		if r < 'a' || r > 'z' {
			return false
		}
	}

	if opts.EngLettersOnly {
		for _, r := range msg {
			if unicode.IsLetter(r) && ((r < 'a' || r > 'z') && (r < 'A' || r > 'Z')) {
				return false
			}
		}
	}

	if opts.NoSpecialSymbols {
		for _, r := range msg {
			if !unicode.In(r, unicode.Letter, unicode.Space, unicode.Number) &&
				!slices.Contains(opts.AllowedSymbols, r) {
				return false
			}
		}
	}

	for _, re := range opts.SensitiveKeywordsRe {
		if re.MatchString(msg) {
			return false
		}
	}

	for _, pattern := range opts.SensitivePatterns {
		if pattern.MatchString(msg) {
			return false
		}
	}

	return true
}

func applyStartsWithLowercaseFixText(msg string) (string, bool) {
	s := msg
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r >= 'A' && r <= 'Z' {
			return string(unicode.ToLower(r)) + msg[len(msg)-len(s)+size:], true
		}
		if r < 'a' || r > 'z' {
			s = s[size:]
			continue
		}

		return msg[len(msg)-len(s):], true
	}

	return "", false
}

func applyEngLettersOnlyFixText(msg string) (string, bool) {
	builder := strings.Builder{}

	for _, r := range msg {
		if unicode.IsLetter(r) && ((r < 'a' || r > 'z') && (r < 'A' || r > 'Z')) {
			continue
		}

		builder.WriteRune(r)
	}

	fixed := builder.String()
	fixed = strings.Join(strings.Fields(fixed), " ")
	if fixed == msg {
		return "", false
	}

	return fixed, true
}

func applyNoSpecialSymbolsFixText(opts *Options, msg string) (string, bool) {
	builder := strings.Builder{}

	for _, r := range msg {
		if !unicode.In(r, unicode.Letter, unicode.Space, unicode.Number) && !slices.Contains(opts.AllowedSymbols, r) {
			continue
		}

		builder.WriteRune(r)
	}

	fixed := builder.String()
	fixed = strings.Join(strings.Fields(fixed), " ")
	if fixed == msg {
		return "", false
	}

	return fixed, true
}

func applyNoSensitiveKeywordsFixText(opts *Options, msg string) (string, bool) {
	fixed := msg

	for _, re := range opts.SensitiveKeywordsRe {
		fixed = re.ReplaceAllString(fixed, "")
	}

	fixed = strings.Join(strings.Fields(fixed), " ")
	if fixed == msg {
		return "", false
	}

	return fixed, true
}

func applyNoSensitivePatternsFixText(opts *Options, msg string) (string, bool) {
	fixed := msg

	for _, pattern := range opts.SensitivePatterns {
		fixed = pattern.ReplaceAllString(fixed, "")
	}

	fixed = strings.Join(strings.Fields(fixed), " ")
	if fixed == msg {
		return "", false
	}

	return fixed, true
}
