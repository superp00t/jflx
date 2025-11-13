package tool

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/superp00t/jflx/scraper"
)

var ScrapeTVShow = cobra.Command{
	Use: "tvshow path",
	Run: run_scrape_tvshow,
}

func init() {
	f := ScrapeTVShow.Flags()
	f.StringP("tmdb-api-key", "k", "", "the TMDB API key")
}

func run_scrape_tvshow(c *cobra.Command, args []string) {
	f := c.Flags()
	var p scraper.ScrapeTVShowParams
	var err error
	p.TMDB_API_Key, err = f.GetString("tmdb-api-key")
	if err != nil {
		log.Fatal(err)
		return
	}
	p.TVShowName = args[0]
	if err = scraper.ScrapeTVShow(&p); err != nil {
		log.Fatal(err)
		return
	}
}
