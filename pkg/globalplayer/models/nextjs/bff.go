package nextjs

type BffPlayableResponse struct {
	Id          string        `json:"id"`
	ContentType string        `json:"contentType"`
	Title       string        `json:"title"`
	Playback    []BffPlayback `json:"playback"`
}

type BffPlayback struct {
	Flags []string `json:"flags"`
	Url   string   `json:"url"`
}
