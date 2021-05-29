package main

import (
	flag "github.com/spf13/pflag"
	"github.com/webview/webview"
	"github.com/zserge/lorca"
	"os"
)

func main() {
	url := flag.StringP("url", "u", "https://www.arsenm.dev", "URL to open in webview")
	debug := flag.BoolP("debug", "d", false, "Enable webview debug mode")
	width := flag.IntP("width", "w", 800, "Width of webview window")
	height := flag.IntP("height", "h", 600, "Height of webview window")
	chrome := flag.Bool("chrome", false, "Use chrome devtools protocol via lorca instead of webview")
	flag.Parse()

	if *chrome {
		// If chrome does not exist
		if lorca.LocateChrome() == "" {
			// Display download prompt
			lorca.PromptDownload()
			// Exit with code 1
			os.Exit(1)
		}
		// Create new lorca UI
		l, _ := lorca.New(*url, "", *width, *height)
		defer l.Close()
		// Wait until window closed
		<-l.Done()
	} else {
		// Create new webview
		w := webview.New(*debug)
		defer w.Destroy()
		// Set title of webview window
		w.SetTitle("WebView SSB")
		// Set window size
		w.SetSize(*width, *height, webview.HintNone)
		// Navigate to specified URL
		w.Navigate(*url)
		// Run app
		w.Run()
	}
}
