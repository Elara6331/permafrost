package main

const DesktopTemplate = `#!/usr/bin/env xdg-open
[Desktop Entry]
Name=%s
Icon=%s
Type=Application
Terminal=false
Exec=webview-permafrost --url %s
Categories=%s;`
