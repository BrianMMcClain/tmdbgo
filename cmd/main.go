package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/brianmmcclain/tmdbgo"
)

func main() {
	id := flag.String("id", "", "Movie ID")
	flag.Parse()

	if *id == "" {
		slog.Error("Must provide -id flag")
		os.Exit(-1)
	}

	tmdbKey := os.Getenv("TMDB_KEY")

	if tmdbKey == "" {
		slog.Error("Must provide API key")
		os.Exit(-1)
	}

	tmdb := tmdbgo.NewTMDB(tmdbKey)
	m := tmdb.GetMovie(*id)

	fmt.Println(m)
}
