package main

import (
	"fmt"
	"log"

	"context"
	"github.com/pkg/browser"
	"net/http"
	"net/url"
)

func GetOAuth(client_id string) (string, string, error) {
	url_params := url.Values{}
	url_params.Add("client_id", client_id)
	url_params.Add("redirect_uri", "http://localhost")
	url_params.Add("response_type", "token")
	url_params.Add("scope", "")

	twitch_url := "https://id.twitch.tv/oauth2/authorize?" + url_params.Encode()

	err := browser.OpenURL(twitch_url)
	if err != nil {
		fmt.Println("Please Visit:", twitch_url)
	}

	m := http.NewServeMux()
	s := http.Server{Addr: "localhost:80", Handler: m}

	access_token := ""

	m.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add check for favicon stuff

		url_params := r.URL.Query()

		if temp, err := url_params["access_token"]; err {
			access_token = temp[0]
		}

		log.Println("Shutting Down HTTP Server")
		s.Shutdown(context.Background())
	})

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add check for favicon stuff

		url_params := r.URL.Query()
		if _, err := url_params["error"]; err { // TODO: Make this return an error in the parent function
			log.Fatalln("Twitch Authentication Error:", url_params["error"][0], url_params["error_description"][0])
		}

		// Twitch sends back access token data specifically in the form of a hash in the URL
		// so that it isn't sent along with the http request. This of course, means that the
		// server never gets it, so I need to load some Javascript to do the very thing Twitch
		// didn't want me to do.
		// Don't you just love security?
		fmt.Fprintln(w, `<script>
			document.location.href = "http://localhost/token?" + document.location.hash.substr(1);
		</script>`)
	})

	log.Println("Starting HTTP Server")
	if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return "", err
	}

	return access_token, nil
}
