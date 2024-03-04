package models

type StationsPageResponse struct {
	PageProps StatsionsPageProps `json:"pageProps"`
}

type StatsionsPageProps struct {
	Feature Feature `json:"feature"`
}

type Feature struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Brands                  []StationBrand   `json:"brands,omitempty"`
	DependsOn               []string         `json:"dependsOn,omitempty"`
	FetchBeforePresentation string           `json:"fetchBeforePresentation,omitempty"`
	Identifier              string           `json:"identifier"`
	PreferredBrands         []PrefferedBrand `json:"preferredBrands,omitempty"`
	Styles                  []any            `json:"styles"`
	Subtitle                any              `json:"subtitle,omitempty"`
	Timestamp               int              `json:"timestamp,omitempty"`
	Title                   string           `json:"title,omitempty"`
	Type                    string           `json:"type"`
	Item                    struct {
		ContentShape string `json:"contentShape"`
		Description  string `json:"description"`
		Flags        []any  `json:"flags"`
		Image        struct {
			Shape string `json:"shape"`
			URL   string `json:"url"`
		} `json:"image"`
		Link struct {
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"link"`
		Metadata         string `json:"metadata"`
		NodeID           string `json:"nodeId"`
		RecommendationID any    `json:"recommendationId"`
		Title            string `json:"title"`
		Video            any    `json:"video"`
	} `json:"item,omitempty"`
	OrnamentColour string `json:"ornament_colour,omitempty"`
	Theme          string `json:"theme,omitempty"`
	Items          []struct {
		Image struct {
			Shape string `json:"shape"`
			URL   string `json:"url"`
		} `json:"image"`
		Link struct {
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"link"`
		Metadata         string `json:"metadata"`
		RecommendationID any    `json:"recommendationId"`
		Subtitle         string `json:"subtitle"`
		Title            string `json:"title"`
	} `json:"items,omitempty"`
}

type StationBrand struct {
	ID                             string          `json:"id"`
	LogoUnstackedWithoutBackground string          `json:"logoUnstackedWithoutBackground"`
	Name                           string          `json:"name"`
	NationalStation                NationalStation `json:"nationalStation"`
	ShowBackgroundImage            string          `json:"showBackgroundImage"`
	Slug                           string          `json:"slug"`
}

type NationalStation struct {
	HeraldID  int    `json:"heraldId"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	StreamURL string `json:"streamUrl"`
	Tagline   string `json:"tagline"`
}

type PrefferedBrand struct {
	ID                             string `json:"id"`
	LogoStackedWithBrandBackground string `json:"logoStackedWithBrandBackground"`
	Name                           string `json:"name"`
	Source                         string `json:"source"`
}
