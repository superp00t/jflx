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

# download code
git clone https://github.com/superp00t/jflx.git
# go into josh flicks directory
cd jflx
# now you should create your server.conf in this directory

# build the server code
go build -o jflx_server github.com/superp00t/jflx/cmd/jflx_server
# run your server (linux)
./jflx_server
# or on windows
jflx_server.exe

```