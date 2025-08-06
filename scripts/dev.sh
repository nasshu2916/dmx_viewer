#!/bin/zsh -e

cd "$(dirname "$0")/.."

npm run --prefix ./frontend watch -- --mode development &
(cd ./backend/ && air) &
wait
