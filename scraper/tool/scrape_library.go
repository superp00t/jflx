package tool

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/superp00t/jflx/config"
	"github.com/superp00t/jflx/scraper"
)

var ScrapeLibrary = cobra.Command{
	Use: "library",
	Run: run_scrape_library,
}

func init() {
	f := ScrapeLibrary.Flags()
	f.StringP("tmdb-api-key", "k", "", "the TMDB API key")
	f.StringP("config", "c", "server.conf", "the server config")
}

func run_scrape_library(c *cobra.Command, args []string) {
	f := c.Flags()
	var sp scraper.ScrapeLibraryParams
	var err error
	sp.TMDB_API_Key, err = f.GetString("tmdb-api-key")
	if err != nil {
		log.Fatal(err)
		return
	}
	var config_location string
	config_location, err = f.GetString("config")
	if err != nil {
		log.Fatal(err)
		return
	}
	var server_config config.ServerConfig
	if err = server_config.Load(config_location); err != nil {
		log.Fatal(err)
		return
	}

	sp.Volumes = server_config.Volumes
	if err = scraper.ScrapeLibrary(&sp); err != nil {
		log.Fatal(err)
		return
	}
}
