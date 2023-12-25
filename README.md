# josh flicks

An example server.conf

```json
{
 "ListenAddress": "0.0.0.0:38088",
 "TMDBScrapeKey": "<your TMDB scraping API key here>",
 "Volumes": [
  {
    "Kinds": "movie",
    "Handle": "Film",
    "Sources": [
      "/var/film1/",
      "/media/external_film_drive/"
    ]
  }
 ]
}
```

```bash

## create server.conf

git clone https://github.com/superp00t/jflx.git
cd jflx
go build github.com/superp00t/jflx/cmd/jflx_server -o jflx_server
./jflx_server

```