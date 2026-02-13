package starts_with_lowercase

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func testsSlogMessages() {
	ctx := context.Background()
	l := slog.Default()

	const (
		slogConstUpper = "SlogConstUpper"
		slogConstDigit = "7slogConst"
		slogConstSpace = " slogConst"
		slogConstEmpty = ""
		slogConstCyr   = "ПриветConst"
		slogConstNoEng = "Привет"
	)

	var slogVar = "VarUpper"

	slog.Debug("valid lower")
	slog.Info("Upper literal")                  // want "log message must start with a lowercase English letter"
	slog.Warn("9digit literal")                 // want "log message must start with a lowercase English letter"
	slog.Error(" leading space literal")        // want "log message must start with a lowercase English letter"
	slog.DebugContext(ctx, "")                  // want "log message must start with a lowercase English letter"
	slog.InfoContext(ctx, "   ")                // want "log message must start with a lowercase English letter"
	slog.WarnContext(ctx, "ПриветUpper")        // want "log message must start with a lowercase English letter"
	slog.ErrorContext(ctx, "Приветlower")       // want "log message must start with a lowercase English letter"
	slog.Log(ctx, slog.LevelInfo, "!symbolMsg") // want "log message must start with a lowercase English letter"
	slog.LogAttrs(ctx, slog.LevelInfo, "A")     // want "log message must start with a lowercase English letter"

	slog.Debug(slogConstUpper)                          // want "log message must start with a lowercase English letter"
	slog.InfoContext(ctx, slogConstDigit)               // want "log message must start with a lowercase English letter"
	slog.WarnContext(ctx, slogConstSpace)               // want "log message must start with a lowercase English letter"
	slog.ErrorContext(ctx, slogConstEmpty)              // want "log message must start with a lowercase English letter"
	slog.Log(ctx, slog.LevelWarn, slogConstCyr)         // want "log message must start with a lowercase English letter"
	slog.LogAttrs(ctx, slog.LevelError, slogConstNoEng) // want "log message must start with a lowercase English letter"

	l.Debug(slogVar)
}

func testsZapMessages() {
	l := &zap.Logger{}
	s := &zap.SugaredLogger{}

	const (
		zapConstUpper = "ZapConstUpper"
		zapConstDigit = "8zapConst"
	)

	var zapVar = "VarUpper"

	l.Debug("ZapUpper") // want "log message must start with a lowercase English letter"
	l.Info("1zapDigit") // want "log message must start with a lowercase English letter"
	l.Warn(" zapSpace") // want "log message must start with a lowercase English letter"
	l.Error("validzap")
	l.DPanic("ПриветZapUpper") // want "log message must start with a lowercase English letter"
	l.Panic("Приветzaplower")  // want "log message must start with a lowercase English letter"
	l.Fatal("...zapSymbol")    // want "log message must start with a lowercase English letter"

	s.Debug("SugaredUpper") // want "log message must start with a lowercase English letter"
	s.Info("2sugar")        // want "log message must start with a lowercase English letter"
	s.Warn(" sugared")      // want "log message must start with a lowercase English letter"
	s.Error("")             // want "log message must start with a lowercase English letter"
	s.DPanic(" ")           // want "log message must start with a lowercase English letter"
	s.Panic("Привет")       // want "log message must start with a lowercase English letter"
	s.Fatal("zgood")

	s.Debugf("FmtUpper") // want "log message must start with a lowercase English letter"
	s.Infof("3fmt")      // want "log message must start with a lowercase English letter"
	s.Warnf("\tfmtTab")  // want "log message must start with a lowercase English letter"
	s.Errorf("fmtok")
	s.DPanicf("Äfmt") // want "log message must start with a lowercase English letter"
	s.Panicf("Z")     // want "log message must start with a lowercase English letter"
	s.Fatalf("a")

	s.Debugw("WrapUpper") // want "log message must start with a lowercase English letter"
	s.Infow("4wrap")      // want "log message must start with a lowercase English letter"
	s.Warnw(" wrap")      // want "log message must start with a lowercase English letter"
	s.Errorw("wrapok")
	s.DPanicw("ПриветWrap") // want "log message must start with a lowercase English letter"
	s.Panicw("Приветwrap")  // want "log message must start with a lowercase English letter"
	s.Fatalw("!wrap")       // want "log message must start with a lowercase English letter"

	s.Debugln("LineUpper") // want "log message must start with a lowercase English letter"
	s.Infoln("5line")      // want "log message must start with a lowercase English letter"
	s.Warnln(" line")      // want "log message must start with a lowercase English letter"
	s.Errorln("lineok")
	s.DPanicln("")       // want "log message must start with a lowercase English letter"
	s.Panicln(" ")       // want "log message must start with a lowercase English letter"
	s.Fatalln("ПриветL") // want "log message must start with a lowercase English letter"

	l.Info(zapConstUpper)  // want "log message must start with a lowercase English letter"
	s.Infof(zapConstDigit) // want "log message must start with a lowercase English letter"

	s.Warn(zapVar)
}
