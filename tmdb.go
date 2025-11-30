package tmdbgo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Movie struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	IMDB         string     `json:"imdb"`
	Overview     string     `json:"overview"`
	Tagline      string     `json:"tagline"`
	Runtime      int        `json:"runtime"`
	ReleaseDate  string     `json:"release_date"`
	PosterPath   string     `json:"poster_path"`
	BackdropPath string     `json:"backdrop_path"`
	Status       string     `json:"status"`
	Language     string     `json:"original_language"`
	Genres       []Genre    `json:"genres"`
	Stream       []Provider `json:"stream"`
	Free         []Provider `json:"free"`
	Buy          []Provider `json:"buy"`
	Rent         []Provider `json:"rent"`
	Ads          []Provider `json:"ads"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Provider struct {
	LogoPath string `json:"logo_path"`
	ID       int    `json:"provider_id"`
	Name     string `json:"provider_name"`
}

type TMDB struct {
	auth string
}

type TMDBSearchResult struct {
	Results []Movie `json:"results"`
}

type ProviderSearchResult struct {
	ID      int                                   `json:"id"`
	Locales map[string]ProviderSearchResultLocale `json:"results"`
}

type ProviderSearchResultLocale struct {
	Link   string     `json:"link"`
	Buy    []Provider `json:"buy"`
	Stream []Provider `json:"flatrate"`
	Rent   []Provider `json:"rent"`
	Ads    []Provider `json:"ads"`
	Free   []Provider `json:"free"`
}

func NewTMDB(auth string) *TMDB {

	tmdb := new(TMDB)
	tmdb.auth = auth

	return tmdb
}

func sendGetRequest(url string, tmdb *TMDB) (int, string) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tmdb.auth)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(body)
}

func (tmdb *TMDB) SearchMovies(title string) []Movie {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?query=%s", url.QueryEscape(title))

	_, sJson := sendGetRequest(url, tmdb)
	results := new(TMDBSearchResult)
	err := json.Unmarshal([]byte(sJson), &results)
	if err != nil {
		log.Fatal(err)
	}

	return results.Results
}

func (tmdb *TMDB) GetMovie(id string) *Movie {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%v?language=en-US", id)

	_, sJson := sendGetRequest(url, tmdb)
	//fmt.Println(sJson)
	movie := new(Movie)
	err := json.Unmarshal([]byte(sJson), &movie)
	if err != nil {
		log.Fatal(err)
	}

	return movie
}

func (tmdb *TMDB) GetWatchProviders(movie *Movie, locale string) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%v/watch/providers?locale=US", movie.ID)

	_, sJson := sendGetRequest(url, tmdb)
	res := new(ProviderSearchResult)
	err := json.Unmarshal([]byte(sJson), &res)
	if err != nil {
		log.Fatal(err)
	}

	movie.Buy = append(movie.Buy, res.Locales[locale].Buy...)
	movie.Stream = append(movie.Stream, res.Locales[locale].Stream...)
	movie.Rent = append(movie.Rent, res.Locales[locale].Rent...)
	movie.Ads = append(movie.Ads, res.Locales[locale].Ads...)
	movie.Free = append(movie.Free, res.Locales[locale].Free...)
}
