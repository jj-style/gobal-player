package models

type CatchupResponse struct {
	PageProps CatchupPageProps `json:"pageProps"`
}

type CatchupPageProps struct {
	CatchupInfo []CatchupInfo `json:"catchupInfo"`
	BrandData   []Brand       `json:"brandData"`
}

type CatchupInfo struct {
	ID       string `json:"id"`
	ImageURL string `json:"imageUrl"`
	Title    string `json:"title"`
}
