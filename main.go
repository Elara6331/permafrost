package main

import (
	"fmt"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	flag "github.com/spf13/pflag"
	"image"
	"os"
	"runtime"
)

const name = "permafrost"

var home string

func init() {
	// If not on linux, fatally log
	if runtime.GOOS != "linux" {
		log.Fatal().Msg("This tool only supports Linux.")
	}
	var err error
	// Get user home directory
	home, err = os.UserHomeDir()
	if err != nil {
		errDisp(false, err, "Error getting user home directory")
	}
}

func main() {
	useGui := flag.BoolP("gui", "g", false, "Use GUI (ignores all other flags)")
	create := flag.BoolP("create", "c", false, "Create new SSB")
	remove := flag.BoolP("remove", "r", false, "Remove an existing SSB")
	ssbName := flag.StringP("name", "n", "", "Name of SSB to create or remove")
	url := flag.StringP("url", "u", "", "URL of new SSB")
	category := flag.StringP("category", "C", "Network", "Category of new SSB")
	iconPath := flag.StringP("icon", "i", "", "Path to icon for new SSB")
	chrome := flag.Bool("chrome", false, "Use chrome via lorca instead of webview")
	flag.ErrHelp = fmt.Errorf("help message for %s", name)
	flag.Parse()

	// If --gui provided
	if *useGui {
		// Start GUI
		initGUI()
	} else {
		// If --create provided
		if *create {
			// Open icon path provided via --icon
			iconFile, err := os.Open(*iconPath)
			if err != nil {
				errDisp(false, err, "Error opening icon file")
			}
			defer iconFile.Close()
			// Decode icon file into image.Image
			icon, _, err := image.Decode(iconFile)
			if err != nil {
				errDisp(false, err, "Error decoding image from file")
			}
			// Attempt to create SSB using flag-provided values
			err = createSSB(*ssbName, *url, *category, *chrome, icon)
			if err != nil {
				errDisp(false, err, "Error creating SSB")
			}
		} else if *remove {
			// Attempt to remove ssb of name provided via --name
			err := removeSSB(*ssbName)
			if err != nil {
				errDisp(false, err, "Error removing SSB")
			}
		} else {
			// Show help screen
			flag.Usage()
			log.Fatal().Msg("Must provide --gui, --create, or --remove")
		}
	}
}

func initGUI() {
	app := fyneApp.New()
	// Create new window with title
	window := app.NewWindow("Webview SSB")

	// Create tab container
	tabs := container.NewAppTabs(
		container.NewTabItem("Create", createTab(window)),
		container.NewTabItem("Remove", removeTab(window)),
	)
	// Put tabs at the top of the window
	tabs.SetTabLocation(container.TabLocationTop)

	window.SetContent(tabs)
	window.ShowAndRun()
}
