package plugin

import (
	"github.com/ddovwll/loglint"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

const Name = "loglint"

type Settings struct {
	StartsWithLowercase *bool     `json:"start-with-lowercase"`
	EngLettersOnly      *bool     `json:"eng-letters"`
	NoSpecialSymbols    *bool     `json:"no-special-symbols"`
	AllowedSymbols      *string   `json:"allowed-symbols"`
	SensitiveKeywords   *[]string `json:"sensitive-keywords"`
	SensitivePatterns   *[]string `json:"sensitive-patterns"`
}

type Plugin struct {
	cfg Settings
}

func New(settings any) (register.LinterPlugin, error) {
	cfg, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}
	return &Plugin{cfg: cfg}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	opts, err := p.createOptions()
	if err != nil {
		return nil, err
	}

	return []*analysis.Analyzer{
		loglint.NewAnalyzer(opts),
	}, nil
}

func (p *Plugin) createOptions() (*loglint.Options, error) {
	builder := loglint.NewOptionsBuilder()

	if p.cfg.StartsWithLowercase != nil {
		builder.SetStartsWithLowercase(*p.cfg.StartsWithLowercase)
	}
	if p.cfg.EngLettersOnly != nil {
		builder.SetEngLettersOnly(*p.cfg.EngLettersOnly)
	}
	if p.cfg.NoSpecialSymbols != nil {
		builder.SetNoSpecialSymbols(*p.cfg.NoSpecialSymbols)
	}
	if p.cfg.AllowedSymbols != nil {
		builder.SetAllowedSymbols(*p.cfg.AllowedSymbols)
	}
	if p.cfg.SensitiveKeywords != nil {
		builder.SetSensitiveKeywords(*p.cfg.SensitiveKeywords)
	}
	if p.cfg.SensitivePatterns != nil {
		builder.SetSensitivePatterns(*p.cfg.SensitivePatterns)
	}

	built, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return built, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func init() {
	register.Plugin(Name, New)
}
