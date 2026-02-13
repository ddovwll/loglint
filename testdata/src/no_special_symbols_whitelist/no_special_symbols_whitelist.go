package no_special_symbols_whitelist

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func testsSlogMessages() {
	ctx := context.Background()
	l := slog.Default()

	const (
		slogConstGood = "const_ok.v1:ready\U0001F600"
		slogConstBad  = "const!bad"
	)

	var slogVar = "var;bad"

	slog.Debug("ok_under_score_1")
	slog.Info("ok.dot.v2")
	slog.Warn("ok:colon:3")
	slog.Error("ok\U0001F600emoji")
	slog.DebugContext(ctx, "mix!bad_under")          // want "found forbidden symbol"
	slog.InfoContext(ctx, "semi;bad.dot")            // want "found forbidden symbol"
	slog.WarnContext(ctx, "dash-bad:colon")          // want "found forbidden symbol"
	slog.ErrorContext(ctx, "slash/bad_x")            // want "found forbidden symbol"
	slog.Log(ctx, slog.LevelInfo, "line\nbad:x")     // want "found forbidden symbol"
	slog.LogAttrs(ctx, slog.LevelInfo, "tab\tbad.y") // want "found forbidden symbol"

	slog.Debug(slogConstGood)
	slog.InfoContext(ctx, slogConstBad) // want "found forbidden symbol"

	l.Debug(slogVar)
	l.Info("ok_logger.value_1")
	l.Warn("hash#bad:warn") // want "found forbidden symbol"
	l.Error("ok:logger\U0001F600x")
	l.DebugContext(ctx, "quest?bad_ok") // want "found forbidden symbol"
	l.InfoContext(ctx, "ok.path:_x")
	l.WarnContext(ctx, "plus+bad.z") // want "found forbidden symbol"
	l.ErrorContext(ctx, "ok\U0001F600:ctx")
	l.Log(ctx, slog.LevelWarn, "at@bad.path") // want "found forbidden symbol"
	l.LogAttrs(ctx, slog.LevelError, "ok_attr.value:9")
}

func testsZapMessages() {
	l := &zap.Logger{}
	s := &zap.SugaredLogger{}

	const (
		zapConstGood = "zap_ok.v1:done\U0001F600"
		zapConstBad  = "zap/bad"
	)

	var zapVar = "zap@var"

	l.Debug("ok_l.debug")
	l.Info("bad!info") // want "found forbidden symbol"
	l.Warn("ok:l_warn")
	l.Error("bad,error") // want "found forbidden symbol"
	l.DPanic("ok\U0001F600panic")
	l.Panic("bad;panic") // want "found forbidden symbol"
	l.Fatal("ok_final:1")

	s.Debug("ok_s.debug")
	s.Info("bad/info") // want "found forbidden symbol"
	s.Warn("ok:warn\U0001F600")
	s.Error("bad\\err") // want "found forbidden symbol"
	s.DPanic("ok_dp.value")
	s.Panic("bad=panic") // want "found forbidden symbol"
	s.Fatal("ok_fatal:2")

	s.Debugf("ok_f.debug")
	s.Infof("bad<fmt") // want "found forbidden symbol"
	s.Warnf("ok:fmt\U0001F600")
	s.Errorf("bad>fmt") // want "found forbidden symbol"
	s.DPanicf("ok_fmt.value")
	s.Panicf("bad*fmt") // want "found forbidden symbol"
	s.Fatalf("ok_fmt:3")

	s.Debugw("ok_w.debug")
	s.Infow("bad^w") // want "found forbidden symbol"
	s.Warnw("ok:w\U0001F600")
	s.Errorw("bad&w") // want "found forbidden symbol"
	s.DPanicw("ok_w.value")
	s.Panicw("bad\"w") // want "found forbidden symbol"
	s.Fatalw("ok_w:4")

	s.Debugln("ok_ln.debug")
	s.Infoln("bad'w") // want "found forbidden symbol"
	s.Warnln("ok:ln\U0001F600")
	s.Errorln("bad(ln") // want "found forbidden symbol"
	s.DPanicln("ok_ln.value")
	s.Panicln("bad)ln") // want "found forbidden symbol"
	s.Fatalln("ok_ln:5")

	l.Info(zapConstGood)
	s.Infof(zapConstBad) // want "found forbidden symbol"

	s.Warn(zapVar)
}
