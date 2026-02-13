package no_sensitive_keywords

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
		slogConstBadIP   = "const ip_address value"
		slogConstBadPet  = "const name_of_the_first_pet value"
		slogConstOK      = "const safe message"
	)

	var slogVar = "var password value"

	slog.Debug("safe slog debug")
	slog.Info("password leaked")                                               // want "found sensitive keyword"
	slog.Warn("client ip_address blocked")                                     // want "found sensitive keyword"
	slog.Error("name_of_the_first_pet required")                               // want "found sensitive keyword"
	slog.DebugContext(ctx, "PASSWORD reset")                                   // want "found sensitive keyword"
	slog.InfoContext(ctx, "Ip_Address missing")                                // want "found sensitive keyword"
	slog.WarnContext(ctx, "Name_Of_The_First_Pet wrong")                       // want "found sensitive keyword"
	slog.ErrorContext(ctx, "password ip_address pair")                         // want "found sensitive keyword"
	slog.Log(ctx, slog.LevelInfo, "before password after")                     // want "found sensitive keyword"
	slog.LogAttrs(ctx, slog.LevelWarn, "name_of_the_first_pet and ip_address") // want "found sensitive keyword"

	slog.Debug(slogConstBadPass)           // want "found sensitive keyword"
	slog.InfoContext(ctx, slogConstBadIP)  // want "found sensitive keyword"
	slog.WarnContext(ctx, slogConstBadPet) // want "found sensitive keyword"
	slog.ErrorContext(ctx, slogConstOK)

	l.Debug(slogVar)
	l.Info("password")               // want "found sensitive keyword"
	l.Warn("ip_address")             // want "found sensitive keyword"
	l.Error("name_of_the_first_pet") // want "found sensitive keyword"
	l.DebugContext(ctx, "safe logger context")
	l.InfoContext(ctx, "password in context")               // want "found sensitive keyword"
	l.WarnContext(ctx, "ip_address in context")             // want "found sensitive keyword"
	l.ErrorContext(ctx, "name_of_the_first_pet in context") // want "found sensitive keyword"
	l.Log(ctx, slog.LevelError, "password and ip_address")  // want "found sensitive keyword"
	l.LogAttrs(ctx, slog.LevelError, "safe attrs")
}

func testsZapMessages() {
	l := &zap.Logger{}
	s := &zap.SugaredLogger{}

	const (
		zapConstBad = "const password in zap"
		zapConstOK  = "const safe zap"
	)

	var zapVar = "var ip_address in zap"

	l.Debug("password in debug") // want "found sensitive keyword"
	l.Info("safe info")
	l.Warn("ip_address warn")              // want "found sensitive keyword"
	l.Error("name_of_the_first_pet error") // want "found sensitive keyword"
	l.DPanic("PASSWORD")                   // want "found sensitive keyword"
	l.Panic("Ip_Address value")            // want "found sensitive keyword"
	l.Fatal("safe fatal")

	s.Debug("name_of_the_first_pet in sugar") // want "found sensitive keyword"
	s.Info("safe sugar")
	s.Warn("password")                 // want "found sensitive keyword"
	s.Error("ip_address and password") // want "found sensitive keyword"
	s.DPanic("safe dpanic")
	s.Panic("Name_Of_The_First_Pet and ip_address") // want "found sensitive keyword"
	s.Fatal("safe fatal sugar")

	s.Debugf("password fmt") // want "found sensitive keyword"
	s.Infof("safe fmt")
	s.Warnf("ip_address fmt")             // want "found sensitive keyword"
	s.Errorf("name_of_the_first_pet fmt") // want "found sensitive keyword"
	s.DPanicf("safe dpanicf")
	s.Panicf("password and ip_address fmt") // want "found sensitive keyword"
	s.Fatalf("safe fatalf")

	s.Debugw("password w") // want "found sensitive keyword"
	s.Infow("safe w")
	s.Warnw("ip_address w")             // want "found sensitive keyword"
	s.Errorw("name_of_the_first_pet w") // want "found sensitive keyword"
	s.DPanicw("safe dpanicw")
	s.Panicw("password and ip_address w") // want "found sensitive keyword"
	s.Fatalw("safe fatalw")

	s.Debugln("password ln") // want "found sensitive keyword"
	s.Infoln("safe ln")
	s.Warnln("ip_address ln")             // want "found sensitive keyword"
	s.Errorln("name_of_the_first_pet ln") // want "found sensitive keyword"
	s.DPanicln("safe dpanicln")
	s.Panicln("password and ip_address ln") // want "found sensitive keyword"
	s.Fatalln("safe fatalln")

	l.Info(zapConstBad) // want "found sensitive keyword"
	s.Infof(zapConstOK)

	s.Warn(zapVar)
}
