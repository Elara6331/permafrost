package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"os"
)

// Set global logger to zerolog
var log = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})

func errDisp(gui bool, err error, msg string, window ...fyne.Window) {
	// If gui is being used
	if gui {
		// Create new container with message label
		content := container.NewVBox(
			widget.NewLabel(msg),
		)
		if err != nil {
			// Add more details dropdown with error label
			content.Add(widget.NewAccordion(
				widget.NewAccordionItem("More Details", widget.NewLabel(err.Error())),
			))
		}
		// Create and show new custom dialog with container
		dialog.NewCustom("Error", "Ok", content, window[0]).Show()
	} else {
		if err != nil {
			log.Warn().Err(err).Msg(msg)
		} else {
			log.Warn().Msg(msg)
		}
	}
}
