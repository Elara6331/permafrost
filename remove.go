package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
	"path/filepath"
)

var rmBtns *fyne.Container

func removeTab(window fyne.Window) *fyne.Container {
	// Use home directory to get various paths
	configDir := filepath.Join(home, ".config", name)
	iconDir := filepath.Join(configDir, "icons")

	// Create directories if they do not exist
	err := makeDirs(configDir, iconDir)
	if err != nil {
		errDisp(true, err, "Error creating required directories", window)
	}

	// Create new wrapping grid for remove buttons
	rmBtns = container.NewGridWrap(fyne.NewSize(125, 75))

	// Refresh remove buttons, adding any existing SSBs
	refreshRmBtns(window, iconDir, rmBtns)

	return rmBtns
}

func refreshRmBtns(window fyne.Window, iconDir string, rmBtns *fyne.Container) {
	// Remove all objects from container
	rmBtns.Objects = []fyne.CanvasObject{}
	// List files in icon directory
	ls, err := os.ReadDir(iconDir)
	if err != nil {
		errDisp(true, err, "Error listing icon directory", window)
	}
	for _, listing := range ls {
		listingName := listing.Name()
		// Get path for SSB icon
		listingPath := filepath.Join(iconDir, listingName)

		// Load icon from path
		img, err := fyne.LoadResourceFromPath(filepath.Join(listingPath, "icon.png"))
		if err != nil {
			errDisp(true, err, "Error loading icon as resource", window)
		}

		// Create new button with icon
		rmBtn := widget.NewButtonWithIcon(listingName, img, func() {
			// Create and show new confirmation dialog
			dialog.NewConfirm(
				"Remove SSB",
				fmt.Sprintf("Are you sure you want to remove %s?", listingName),
				func(ok bool) {
					if ok {
						// Attempt to remove SSB
						err = removeSSB(listingName)
						if err != nil {
							errDisp(true, err, "Error removing SSB", window)
						}
						refreshRmBtns(window, iconDir, rmBtns)
					}
				},
				window,
			).Show()
		})
		// Add button to container
		rmBtns.Objects = append(rmBtns.Objects, rmBtn)
	}
	// Refresh container (update changes)
	rmBtns.Refresh()
}

func removeSSB(ssbName string) error {
	// Use home directory to get various paths
	configDir := filepath.Join(home, ".config", name)
	iconDir := filepath.Join(configDir, "icons", ssbName)
	desktopDir := filepath.Join(home, ".local", "share", "applications")

	// Remove icon directory
	err := os.RemoveAll(iconDir)
	if err != nil {
		return err
	}

	// Remove desktop file
	err = os.Remove(filepath.Join(desktopDir, ssbName+".desktop"))
	if err != nil {
		return err
	}
	return nil
}
