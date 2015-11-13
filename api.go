package main

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"fmt"
	"os/exec"
)

// The one and only martini instance.
var m *martini.Martini

const (
	FORM_INTENSITY = "intensity"
	FORM_INTERVAL = "interval"
)

func init() {
	m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(MapEncoder)
	// Setup routes
	r := martini.NewRouter()
	r.Get("/", func() string {	
		fmt.Println("Going to launch program")
		// path, err := exec.LookPath("fswebcam")
		// if err != nil {
		//	log.Fatal("installing fswebcam is in your future")
		//}
		// fmt.Printf("fswebcam is available at %s\n", path)		
		
		cmd := exec.Command("fswebcam", "shot.jpeg")
		err1 := cmd.Run()
		if err1 != nil {
			log.Fatal(err1)
		}
  		return "Hello, webcam" // HTTP 200 : "hello world"
	})
	// r.Post(`/:version/clear`, GlowOff)
	// r.Post(`/:version/on`, TurnAllOn)
	// r.Post(`/:version/flare`, Flare)
	// r.Post(`/:version/colors/:id`, SetGlowColor)
	// r.Post(`/:version/legs/:id`, SetLeg)
	// r.Post(`/:version/legs/:id/colors/:colid`, SetLed)
	// r.Post(`/:version/fan/:num`, SetFan)
	// r.Post(`/:version/fade/:num`, SetFade)

	// Add the router action
	m.Action(r.Handle)
}

// The regex to check for the requested format (allows an optional trailing
// slash).
var rxExt = regexp.MustCompile(`(\.(?:text|json))\/?$`)

// MapEncoder intercepts the request's URL, detects the requested format,
// and injects the correct encoder dependency for this request. It rewrites
// the URL to remove the format extension, so that routes can be defined
// without it.
func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	// Get the format extension
	matches := rxExt.FindStringSubmatch(r.URL.Path)
	ft := ".json"
	if len(matches) > 1 {
		// Rewrite the URL without the format extension
		l := len(r.URL.Path) - len(matches[1])
		if strings.HasSuffix(r.URL.Path, "/") {
			l--
		}
		r.URL.Path = r.URL.Path[:l]
		ft = matches[1]
	}
	// Inject the requested encoder
	switch ft {
	case ".text":
		c.MapTo(textEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	default:
		c.MapTo(jsonEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}

func main() {
	go func() {
		TurnAllOn()
		time.Sleep(3 * time.Second)
		GlowOff()
	}()
	if err := http.ListenAndServe(":9999", m); err != nil {
		log.Fatal(err)
	}
}
