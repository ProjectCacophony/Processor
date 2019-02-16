package lastfm

type Config struct {
	Key    string `envconfig:"LASTFM_KEY"`
	Secret string `envconfig:"LASTFM_SECRET"`
}
