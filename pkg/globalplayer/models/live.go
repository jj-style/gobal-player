package models

type LiveResponse struct {
	PageProps LivePageProps `json:"pageProps"`
}

type LivePageProps struct {
	Brands []Brand `json:"brands"`
}
