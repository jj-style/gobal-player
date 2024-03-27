package nextjs

import "time"

type CatchupShowResponse struct {
	PageProps CatchupShowPageProps `json:"pageProps"`
}

type CatchupShowPageProps struct {
	Station     CatchupShowStationBrand `json:"station"`
	CatchupInfo CatchupInfoDetails      `json:"catchupInfo"`
}

type CatchupInfoDetails struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
	Episodes    []Episode `json:"episodes"`
}

type Episode struct {
	Availability   string    `json:"availability"`
	AvailableUntil time.Time `json:"availableUntil"`
	Description    string    `json:"description"`
	Duration       string    `json:"duration"`
	ID             string    `json:"id"`
	ImageURL       string    `json:"imageUrl"`
	StartDate      time.Time `json:"startDate"`
	StreamURL      string    `json:"streamUrl"`
	Title          string    `json:"title"`
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
