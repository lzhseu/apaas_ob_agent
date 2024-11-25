package logger

func WithFormat(format string) OptionFunc {
	return func(cfg *Config) {
		cfg.format = format
	}
}

func WithLevel(level string) OptionFunc {
	return func(cfg *Config) {
		cfg.level = level
	}
}

func WithAddSource(addSource bool) OptionFunc {
	return func(cfg *Config) {
		cfg.addSource = addSource
	}
}
