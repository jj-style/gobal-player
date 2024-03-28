package nextjs

type Brand struct {
	BrandID             string `json:"brandId"`
	BrandLogo           string `json:"brandLogo"`
	BrandName           string `json:"brandName"`
	BrandSlug           string `json:"brandSlug"`
	Gduid               string `json:"gduid"`
	HeraldID            int    `json:"heraldId"`
	ID                  string `json:"id"`
	LegacyStationPrefix string `json:"legacyStationPrefix"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	StreamURL           string `json:"streamUrl"`
	Tagline             string `json:"tagline"`
}
