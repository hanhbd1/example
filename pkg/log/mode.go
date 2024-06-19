package log

type Mode uint8

const (
	ModeUnknown Mode = iota
	ModeDevelopment
	ModeProduction
)

func (m Mode) String() string {
	return [...]string{
		"unknown", "development", "production",
	}[m]
}

func GetModeFromString(mode string) (m Mode) {
	switch mode {
	case "local":
		m = ModeDevelopment
	case "development", "dev",
		"staging", "stg",
		"production", "prod":
		m = ModeProduction
	default:
		m = ModeUnknown
	}
	if mode == "" {
		m = ModeDevelopment
	}
	return
}
