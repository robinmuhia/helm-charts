package domain

// ImageDetails represents a base docker image
type ImageDetails struct {
	Image  string `json:"image"`
	Size   int64  `json:"size"`
	Layers int    `json:"layers"`
}
