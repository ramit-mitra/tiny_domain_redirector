#!/bin/sh

OUTPUT="build/redirector"

gum style \
	--foreground "#A076F9" --border-foreground "#6528F7" --border double \
	--align center --width 50 --margin "0 0" --padding "2 4" \
	"Redirector build script"

gum style \
	--foreground "#00DFA2" \
	--align left --width 50 --margin "0 0" --padding "1 1" \
	"Clearing old build directory..."
rm -rf $OUTPUT
# also remove app.log
rm -rf app.log

gum style \
	--foreground "#00DFA2" \
	--align left --width 50 --margin "0 0" --padding "1 1" \
	"Building..."
# building for linux x86_64
env GOOS=linux GOARCH=amd64 go build -o $OUTPUT
chmod +x $OUTPUT

gum style \
	--foreground "#2CD3E1" \
	--align left --width 50 --margin "0 0" --padding "1 1" \
	"Done. Build output: $OUTPUT"