package no_sensitive_patterns

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

func testsSlogMessages() {
	ctx := context.Background()
	l := slog.Default()

	const (
		slogConstBadPass = "const password value"
		slogConstBadTok  = "const token value"
		slogConstBadIP   = "const ip_address value"
		slogConstBadCard = "const card 1111-2222-3333-4444 value"
		slogConstOK      = "const mytoken and password123 are safe"
	)

	var slogVar = "var secret value"

	slog.Debug("safe slog debug")
	slog.Info("password leaked")                              // want "found sensitive pattern"
	slog.Warn("client TOKEN blocked")                         // want "found sensitive pattern"
	slog.Error("my secret value")                             // want "found sensitive pattern"
	slog.DebugContext(ctx, "ip_address=1.2.3.4")              // want "found sensitive pattern"
	slog.InfoContext(ctx, "password token secret ip_address") // want "found sensitive pattern"
	slog.WarnContext(ctx, "mytoken and password123 are safe")
	slog.ErrorContext(ctx, "prefix_token suffix")
	slog.Log(ctx, slog.LevelInfo, "before password after")            // want "found sensitive pattern"
	slog.LogAttrs(ctx, slog.LevelWarn, "token and ip_address")        // want "found sensitive pattern"
	slog.Log(ctx, slog.LevelDebug, "card 1234-5678-9012-3456 leaked") // want "found sensitive pattern"

	slog.Debug(slogConstBadPass)                         // want "found sensitive pattern"
	slog.InfoContext(ctx, slogConstBadTok)               // want "found sensitive pattern"
	slog.WarnContext(ctx, slogConstBadIP)                // want "found sensitive pattern"
	slog.LogAttrs(ctx, slog.LevelInfo, slogConstBadCard) // want "found sensitive pattern"
	slog.ErrorContext(ctx, slogConstOK)

	l.Debug(slogVar)
	l.Info("password")                // want "found sensitive pattern"
	l.Warn("token")                   // want "found sensitive pattern"
	l.Error("secret")                 // want "found sensitive pattern"
	l.DebugContext(ctx, "ip_address") // want "found sensitive pattern"
	l.InfoContext(ctx, "safe logger context")
	l.WarnContext(ctx, "mytoken_value")
	l.ErrorContext(ctx, "token value")              // want "found sensitive pattern"
	l.Log(ctx, slog.LevelError, "secret and token") // want "found sensitive pattern"
	l.LogAttrs(ctx, slog.LevelError, "safe attrs")
	l.Log(ctx, slog.LevelWarn, "masked card 1111-2222-3333-4444") // want "found sensitive pattern"
}

func testsZapMessages() {
	l := &zap.Logger{}
	s := &zap.SugaredLogger{}

	const (
		zapConstBad     = "const password and token in zap"
		zapConstBadCard = "const card 2222-3333-4444-5555 in zap"
		zapConstOK      = "const password123 and mytoken safe"
	)

	var zapVar = "var ip_address in zap"

	l.Debug("password in debug") // want "found sensitive pattern"
	l.Info("safe info")
	l.Warn("token warn")         // want "found sensitive pattern"
	l.Error("secret error")      // want "found sensitive pattern"
	l.DPanic("IP_ADDRESS value") // want "found sensitive pattern"
	l.Panic("safe panic")
	l.Fatal("ip_address_token")
	l.Warn("card 9999-8888-7777-6666") // want "found sensitive pattern"

	s.Debug("secret in sugar") // want "found sensitive pattern"
	s.Info("safe sugar")
	s.Warn("password")          // want "found sensitive pattern"
	s.Error("token and secret") // want "found sensitive pattern"
	s.DPanic("safe dpanic")
	s.Panic("ip_address and password") // want "found sensitive pattern"
	s.Fatal("safe fatal sugar")

	s.Debugf("password fmt") // want "found sensitive pattern"
	s.Infof("safe fmt")
	s.Warnf("token fmt")   // want "found sensitive pattern"
	s.Errorf("secret fmt") // want "found sensitive pattern"
	s.DPanicf("safe dpanicf")
	s.Panicf("ip_address and token fmt") // want "found sensitive pattern"
	s.Fatalf("safe fatalf")

	s.Debugw("password w") // want "found sensitive pattern"
	s.Infow("safe w")
	s.Warnw("token w")   // want "found sensitive pattern"
	s.Errorw("secret w") // want "found sensitive pattern"
	s.DPanicw("safe dpanicw")
	s.Panicw("ip_address and token w") // want "found sensitive pattern"
	s.Fatalw("safe fatalw")

	s.Debugln("password ln") // want "found sensitive pattern"
	s.Infoln("safe ln")
	s.Warnln("token ln")   // want "found sensitive pattern"
	s.Errorln("secret ln") // want "found sensitive pattern"
	s.DPanicln("safe dpanicln")
	s.Panicln("ip_address and token ln") // want "found sensitive pattern"
	s.Fatalln("safe fatalln")
	s.Errorf("billing 1234-5678-9012-3456 failed") // want "found sensitive pattern"

	l.Info(zapConstBad)     // want "found sensitive pattern"
	s.Warn(zapConstBadCard) // want "found sensitive pattern"
	s.Infof(zapConstOK)

	s.Warn(zapVar)
}
