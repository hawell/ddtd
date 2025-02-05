package server

type Config struct {
	BindAddress  string `env:"HTTP_ADDRESS"`
	ReadTimeout  int    `env:"HTTP_READ_TIMEOUT" env-default:"3"`
	WriteTimeout int    `env:"HTTP_WRITE_TIMEOUT" env-default:"3"`
	MaxBodyBytes int64  `env:"HTTP_MAX_BODY_SIZE" env-default:"1000000"`
	PanelRoot    string `env:"HTTP_PANEL_ROOT"`
}
