package no_special_symbols

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func testsSlogMessages() {
	ctx := context.Background()
	l := slog.Default()

	const (
		slogConstBadPunct = "const!bad"
		slogConstBadSlash = "const/bad"
		slogConstGood     = "const good 123"
	)

	var slogVar = "var!bad"

	slog.Debug("letters digits 123")
	slog.Info("has!bang")                          // want "found forbidden symbol"
	slog.Warn("comma,here")                        // want "found forbidden symbol"
	slog.Error("path/segment")                     // want "found forbidden symbol"
	slog.DebugContext(ctx, "tab\tvalue")           // want "found forbidden symbol"
	slog.InfoContext(ctx, "line\nbreak")           // want "found forbidden symbol"
	slog.WarnContext(ctx, "under_score")           // want "found forbidden symbol"
	slog.ErrorContext(ctx, "dot.value")            // want "found forbidden symbol"
	slog.Log(ctx, slog.LevelInfo, "colon:value")   // want "found forbidden symbol"
	slog.LogAttrs(ctx, slog.LevelInfo, "at@value") // want "found forbidden symbol"

	slog.Debug(slogConstBadPunct)            // want "found forbidden symbol"
	slog.InfoContext(ctx, slogConstBadSlash) // want "found forbidden symbol"
	slog.WarnContext(ctx, slogConstGood)

	l.Debug(slogVar)
	l.Info("hash#mark")                         // want "found forbidden symbol"
	l.Warn("equal=value")                       // want "found forbidden symbol"
	l.Error("plus+value")                       // want "found forbidden symbol"
	l.DebugContext(ctx, "emoji\U0001F642")      // want "found forbidden symbol"
	l.InfoContext(ctx, "quote\"value")          // want "found forbidden symbol"
	l.WarnContext(ctx, "apost'value")           // want "found forbidden symbol"
	l.ErrorContext(ctx, "bracket[value")        // want "found forbidden symbol"
	l.Log(ctx, slog.LevelWarn, "paren(value")   // want "found forbidden symbol"
	l.LogAttrs(ctx, slog.LevelError, "brace{v") // want "found forbidden symbol"
}

func testsZapMessages() {
	l := &zap.Logger{}
	s := &zap.SugaredLogger{}

	const (
		zapConstBad = "zap!bad"
		zapConstOK  = "zap const ok"
	)

	var zapVar = "zap@var"

	l.Debug("bang!") // want "found forbidden symbol"
	l.Info("ok info")
	l.Warn("slash/") // want "found forbidden symbol"
	l.Error("ok error")
	l.DPanic("dash-") // want "found forbidden symbol"
	l.Panic("ok panic")
	l.Fatal("pipe|") // want "found forbidden symbol"

	s.Debug("paren)") // want "found forbidden symbol"
	s.Info("ok sugar")
	s.Warn("less<") // want "found forbidden symbol"
	s.Error("ok sugar err")
	s.DPanic("greater>") // want "found forbidden symbol"
	s.Panic("ok sugar panic")
	s.Fatal("caret^") // want "found forbidden symbol"

	s.Debugf("star*") // want "found forbidden symbol"
	s.Infof("ok fmt")
	s.Warnf("eq=") // want "found forbidden symbol"
	s.Errorf("ok errf")
	s.DPanicf("back\\") // want "found forbidden symbol"
	s.Panicf("ok panicf")
	s.Fatalf("semi;") // want "found forbidden symbol"

	s.Debugw("amp&") // want "found forbidden symbol"
	s.Infow("ok w")
	s.Warnw("quest?") // want "found forbidden symbol"
	s.Errorw("ok w err")
	s.DPanicw("dq\"") // want "found forbidden symbol"
	s.Panicw("ok wp")
	s.Fatalw("sq'") // want "found forbidden symbol"

	s.Debugln("nl\nx") // want "found forbidden symbol"
	s.Infoln("ok ln")
	s.Warnln("tab\tx") // want "found forbidden symbol"
	s.Errorln("ok ln err")
	s.DPanicln("emoji\U0001F600") // want "found forbidden symbol"
	s.Panicln("ok lnp")
	s.Fatalln("ok lnf")

	l.Info(zapConstBad) // want "found forbidden symbol"
	s.Infof(zapConstOK)

	s.Warn(zapVar)
}
