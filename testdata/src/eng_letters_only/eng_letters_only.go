package eng_letters_only

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func testsSlogMessages() {
	ctx := context.Background()
	l := slog.Default()

	const (
		slogConstBadRu   = "const Ошибка"
		slogConstBadMix  = "constΔmix"
		slogConstGood    = "const english"
		slogConstOnlyBad = "Привет"
	)

	var slogVar = "var Ошибка"

	slog.Debug("plain english")
	slog.Info("error ошибка")            // want "log message must contain only English letters"
	slog.Warn("alpha Ελληνικά")          // want "log message must contain only English letters"
	slog.Error("mixОшибка")              // want "log message must contain only English letters"
	slog.DebugContext(ctx, "café issue") // want "log message must contain only English letters"
	slog.InfoContext(ctx, "a\tß\tb")     // want "log message must contain only English letters"
	slog.WarnContext(ctx, "Привет")      // want "log message must contain only English letters"
	slog.ErrorContext(ctx, "emoji 😀 ok")
	slog.Log(ctx, slog.LevelInfo, "Лог once")     // want "log message must contain only English letters"
	slog.LogAttrs(ctx, slog.LevelInfo, "Ключ: t") // want "log message must contain only English letters"

	slog.Debug(slogConstBadRu)             // want "log message must contain only English letters"
	slog.InfoContext(ctx, slogConstBadMix) // want "log message must contain only English letters"
	slog.WarnContext(ctx, slogConstGood)
	slog.ErrorContext(ctx, slogConstOnlyBad) // want "log message must contain only English letters"

	l.Debug(slogVar)
}

func testsZapMessages() {
	l := &zap.Logger{}
	s := &zap.SugaredLogger{}

	const (
		zapConstBad = "zap Ошибка"
		zapConstOk  = "zap english"
	)

	var zapVar = "zap var Ошибка"

	l.Debug("Ошибка zap") // want "log message must contain only English letters"
	l.Info("zapОшибка")   // want "log message must contain only English letters"
	l.Warn("warn ok")
	l.Error("ßstart")  // want "log message must contain only English letters"
	l.DPanic("Δpanic") // want "log message must contain only English letters"
	l.Panic("Русский") // want "log message must contain only English letters"
	l.Fatal("final ok")

	s.Debug("ñsugar") // want "log message must contain only English letters"
	s.Info("sugar ok")
	s.Warn("warn Ελληνικά") // want "log message must contain only English letters"
	s.Error("error ok")
	s.DPanic("panic Привет") // want "log message must contain only English letters"
	s.Panic("panic ok")
	s.Fatal("fatal Рус") // want "log message must contain only English letters"

	s.Debugf("ßformat %s") // want "log message must contain only English letters"
	s.Infof("fmt ñ value") // want "log message must contain only English letters"
	s.Warnf("fmt ok")
	s.Errorf("éfmt")  // want "log message must contain only English letters"
	s.DPanicf("δfmt") // want "log message must contain only English letters"
	s.Panicf("panicf ok")
	s.Fatalf("Рус fmt") // want "log message must contain only English letters"

	s.Debugw("w Ελληνικά v") // want "log message must contain only English letters"
	s.Infow("w ok")
	s.Warnw("Рус w") // want "log message must contain only English letters"
	s.Errorw("w ok 2")
	s.DPanicw("wßw") // want "log message must contain only English letters"
	s.Panicw("w ok 3")
	s.Fatalw("終w") // want "log message must contain only English letters"

	s.Debugln("ln Привет test") // want "log message must contain only English letters"
	s.Infoln("ln ok")
	s.Warnln("Ωln") // want "log message must contain only English letters"
	s.Errorln("ln ok 2")
	s.DPanicln("中文") // want "log message must contain only English letters"
	s.Panicln("ln ok 3")
	s.Fatalln("끝ln") // want "log message must contain only English letters"

	l.Info(zapConstBad) // want "log message must contain only English letters"
	s.Infof(zapConstOk)

	s.Warn(zapVar)
}
