# jflx

An example server.conf:

- `"ListenAddress"`
  > the IP address to bind to.
- `"TMDBScrapeKey"`
  > You need an API key for the Movie Database to use metadata scraping.

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

# download code
git clone https://github.com/superp00t/jflx.git
# go into jflx directory
cd jflx
# now you should create your server.conf in this directory

# build the server code
go build github.com/superp00t/jflx/cmd/jflx_server
# run your server (linux)
./jflx_server
# or on windows
jflx_server.exe

```