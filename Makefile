all:
	go build
	go build ./cmd/webview-permafrost

install:
	install -Dm755 permafrost $(PREFIX)/usr/bin/permafrost
	install -Dm755 webview-permafrost $(PREFIX)/usr/bin/webview-permafrost
	install -Dm644 permafrost.desktop $(PREFIX)/usr/share/applications/permafrost.desktop
	install -Dm644 permafrost.png $(PREFIX)/usr/share/pixmaps/permafrost.png
