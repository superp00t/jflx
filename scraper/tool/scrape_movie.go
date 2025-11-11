package tool

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/superp00t/jflx/scraper"
)

var ScrapeMovie = cobra.Command{
	Use: "movie path",
	Run: run_scrape_movie,
}

func init() {
	f := ScrapeMovie.Flags()
	f.StringP("tmdb-api-key", "k", "", "the TMDB API key")
}

func run_scrape_movie(c *cobra.Command, args []string) {
	f := c.Flags()
	var sp scraper.ScrapeMovieParams
	var err error
	sp.TMDB_API_Key, err = f.GetString("tmdb-api-key")
	if err != nil {
		log.Fatal(err)
		return
	}
	sp.MovieName = args[0]
	if err = scraper.ScrapeMovie(&sp); err != nil {
		log.Fatal(err)
		return
	}
}
