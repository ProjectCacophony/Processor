package stocks

type Config struct {
	IEXAPISecret string `envconfig:"IEXCLOUD_API_SECRET"`
}
