package nextjs

import "time"

type CatchupShowResponse struct {
	PageProps CatchupShowPageProps `json:"pageProps"`
}

type CatchupShowPageProps struct {
	Station     CatchupShowStationBrand `json:"station"`
	CatchupInfo CatchupInfoDetails      `json:"catchupInfo"`
	ID          string                  `json:"id"`
}

type CatchupInfoDetails struct {
	Title           string `json:"title"`
	CatchupMetadata `json:"metadata"`

	Blocks []CatchupBlock `json:"blocks"`
	// Episodes []Episode `json:"episodes"`

}

type CatchupMetadata struct {
	Description string `json:"description"`
	Author      string `json:"author"`
	Image       struct {
		Url string `json:"url"`
	} `json:"image"`
}

type CatchupBlock struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
	Image      struct {
		Url string `json:"url"`
	} `json:"image"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Items       []Episode `json:"items"`
}

type Episode struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Image       struct {
		Url string `json:"url"`
	} `json:"image"`
	Title   string         `json:"title"`
	Content EpisodeContent `json:"content"`
}

type EpisodeContent struct {
	Type            string    `json:"type"`
	Duration        string    `json:"duration"`
	DurationSeconds int       `json:"durationSeconds"`
	Published       time.Time `json:"published"`
	Expiry          time.Time `json:"Expiry"`
}

type CatchupShowStationBrand struct {
	BrandID             string `json:"brandId"`
	BrandLogo           string `json:"brandLogo"`
	BrandName           string `json:"brandName"`
	BrandSlug           string `json:"brandSlug"`
	Gduid               string `json:"gduid"`
	HeraldID            string `json:"heraldId"`
	ID                  string `json:"id"`
	LegacyStationPrefix string `json:"legacyStationPrefix"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	StreamURL           string `json:"streamUrl"`
	Tagline             string `json:"tagline"`
}
