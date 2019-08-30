package weather

const (
	geocodeEndpoint   = "https://maps.googleapis.com/maps/api/geocode/json?language=en&key=%s&address=%s"
	darkskyForcastURL = "https://darksky.net/forecast/%s,%s/si24"
	darkSkyHexColor   = "#2B86F3"
)

type Config struct {
	GoogleMapsKey string `envconfig:"GOOGLE_MAPS_KEY"`
	DarkSkyKey    string `envconfig:"DARK_SKY_KEY"`
}
