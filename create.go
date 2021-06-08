package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/mat/besticon/ico"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//go:embed permafrost.png
var defLogo []byte

func createTab(window fyne.Window) *fyne.Container {
	// Create entry field for name of SSB
	ssbName := widget.NewEntry()
	ssbName.SetPlaceHolder("App Name")

	// Create entry field for URL of SSB
	ssbURL := widget.NewEntry()
	ssbURL.SetPlaceHolder("App URL")

	// Create dropdown menu for category in desktop file
	category := widget.NewSelectEntry([]string{
		"AudioVideo",
		"Audio",
		"Video",
		"Development",
		"Education",
		"Game",
		"Graphics",
		"Network",
		"Office",
		"Settings",
		"System",
		"Utility",
	})
	category.PlaceHolder = "Category"

	// Get default logo and decode as png
	img, err := png.Decode(bytes.NewReader(defLogo))
	if err != nil {
		errDisp(true, err, "Error decoding default logo", window)
	}
	// Get canvas image from png
	defaultIcon := canvas.NewImageFromImage(img)
	// Create new container for icon with placeholder line
	iconContainer := container.NewMax(canvas.NewLine(color.Black))
	// Set default icon in container
	setIcon(defaultIcon, iconContainer)

	// Create image selection dialog
	selectImg := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		// Close file at end of function
		defer file.Close()
		if err != nil {
			errDisp(true, err, "Error opening file", window)
		}
		// If no file selected, stop further execution of function
		if file == nil {
			return
		}
		// Get image from file reader
		icon := canvas.NewImageFromReader(file, file.URI().Name())
		setIcon(icon, iconContainer)
	}, window)
	// Create filter constrained to images
	selectImg.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))

	// Create button to use favicon as icon
	faviconBtn := widget.NewButton("Use favicon", func() {
		// Attempt to parse URL
		uri, err := url.ParseRequestURI(ssbURL.Text)
		if err != nil {
			errDisp(true, err, "Error parsing URL. Note that the scheme (https://) is required.", window)
			return
		}
		// Attempt to get favicon using DuckDuckGo API
		res, err := http.Get(fmt.Sprintf("https://external-content.duckduckgo.com/ip3/%s.ico", uri.Host))
		if err != nil {
			errDisp(true, err, "Error getting favicon via DuckDuckGo API", window)
			return
		}
		defer res.Body.Close()
		// Attempt to decode data as ico file
		favicon, err := ico.Decode(res.Body)
		if err != nil {
			errDisp(true, err, "Error decoding ico file", window)
			return
		}
		// Get new image from decoded data
		icon := canvas.NewImageFromImage(favicon)
		setIcon(icon, iconContainer)
	})

	// Create vertical container
	col := container.NewVBox(
		category,
		widget.NewButton("Select icon", selectImg.Show),
		faviconBtn,
	)

	// Use home directory to get icon path
	iconDir := filepath.Join(home, ".config", name, "icons")

	useChrome := widget.NewCheck("Use Chrome (not isolated)", nil)

	// Create new button that creates SSB
	createBtn := widget.NewButton("Create", func() {
		// Attempt to create SSB
		err := createSSB(ssbName.Text, ssbURL.Text, category.Text, useChrome.Checked, iconContainer.Objects[0].(*canvas.Image).Image)
		if err != nil {
			errDisp(true, err, "Error creating SSB", window)
		}
		refreshRmBtns(window, iconDir, rmBtns)
		ssbName.SetText("")
		ssbURL.SetText("")
		category.SetText("")
		useChrome.SetChecked(false)
		setIcon(defaultIcon, iconContainer)
	})

	// Create new vertical container
	content := container.New(layout.NewVBoxLayout(),
		ssbName,
		ssbURL,
		// Create dual-column container to house icon and fields
		container.NewGridWithColumns(2,
			container.NewVBox(
				container.NewCenter(widget.NewLabel("Icon")),
				iconContainer,
				useChrome,
			),
			col,
		),
		// Add expanding spacer
		layout.NewSpacer(),
		createBtn,
	)

	// Return base container
	return content
}

func setIcon(img *canvas.Image, logoLayout *fyne.Container) {
	// Expand image keeping aspect ratio
	img.FillMode = canvas.ImageFillContain
	// Set minimum image size to 50x50
	img.SetMinSize(fyne.NewSize(50, 50))
	// Replace image in layout
	logoLayout.Objects[0] = img
	// Refresh layout (update changes)
	logoLayout.Refresh()
}

func createSSB(ssbName, uri, category string, useChrome bool, icon image.Image) error {
	// Parse provided URL for validity
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	// Use home directory to get various paths
	configDir := filepath.Join(home, ".config", name)
	iconDir := filepath.Join(configDir, "icons", ssbName)
	desktopDir := filepath.Join(home, ".local", "share", "applications")

	// Create paths if nonexistent
	err = makeDirs(configDir, iconDir, desktopDir)
	if err != nil {
		return err
	}

	// Get paths to resources
	iconPath := filepath.Join(iconDir, "icon.png")
	desktopPath := filepath.Join(desktopDir, ssbName+".desktop")

	// Create icon file
	iconFile, err := os.Create(iconPath)
	if err != nil {
		return err
	}
	// Encode icon as png, writing to file
	err = png.Encode(iconFile, icon)
	if err != nil {
		return err
	}
	iconFile.Close()

	if useChrome {
		uri += " --chrome"
	}

	// Expand desktop file template with provided data
	desktopStr := fmt.Sprintf(DesktopTemplate,
		ssbName,
		iconPath,
		uri,
		category,
	)

	// Create new executable desktop file
	desktopFile, err := os.OpenFile(desktopPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	// Copy expanded desktop file template to file
	_, err = io.Copy(desktopFile, strings.NewReader(desktopStr))
	if err != nil {
		return err
	}
	desktopFile.Close()

	return nil
}

// Make all directories provided if they do not exist
func makeDirs(dirs ...string) error {
	// For each directory
	for _, dir := range dirs {
		// Create directory and parents if required
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
