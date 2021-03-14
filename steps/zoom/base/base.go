package base

const (
	ZoomAPIV2 = "https://api.zoom.us/v2"
)

type Args struct {
	ZoomToken string `env:"ZOOM_TOKEN"`
}
