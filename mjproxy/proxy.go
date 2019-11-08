package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	config "github.com/spf13/viper"
	"gocv.io/x/gocv"

	"github.com/snowzach/mjpeg"
)

type Stream struct {
	Url  string `json:"url" yaml:"url"`
	Name string `json:"name", yaml:"name"`
}

func main() {

	// Flags
	configFile := flag.String("c", "mjproxy.yaml", "config file")
	flag.Parse()

	// Config
	config.SetConfigFile(*configFile)
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Could not read config file: %v\n", err)
	}

	// Get the streams config
	var streams []Stream
	err = config.UnmarshalKey("streams", &streams)
	if err != nil {
		log.Fatalf("error parsing stream: %+v\n", err)
	}

	for _, streamConfig := range streams {

		stream := mjpeg.NewStream(50 * time.Millisecond)

		go func() {

			// Open capture
			capture, err := gocv.OpenVideoCapture(streamConfig.Url)
			if err != nil {
				fmt.Printf("Error opening capture url: %v: %v\n", streamConfig.Url, err)
				return
			}
			defer capture.Close()

			// create the mjpeg stream
			img := gocv.NewMat()
			defer img.Close()

			for {
				// Read an image
				if ok := capture.Read(&img); !ok {
					return
				}
				if img.Empty() {
					continue
				}

				// re-encode with boxes
				data, err := gocv.IMEncode(".jpg", img)
				if err != nil {
					continue
				}

				stream.UpdateJPEG(data)
			}
		}()

		if username := config.GetString("server.username"); username != "" {
			http.HandleFunc("/"+streamConfig.Name, BasicAuth(stream.ServeHTTP, username, config.GetString("server.password")))
		} else {
			http.Handle("/"+streamConfig.Name, stream)
		}
		log.Printf("Added stream '/%s' from %s\n", streamConfig.Name, streamConfig.Url)

	}

	// start http server
	log.Printf("Server listening on :%d\n", config.GetInt("server.port"))
	log.Fatal(http.ListenAndServe(":"+config.GetString("server.port"), nil))

}

func BasicAuth(handler http.HandlerFunc, username, password string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="mjproxy"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}
