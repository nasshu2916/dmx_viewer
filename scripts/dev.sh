#!/bin/zsh -e

cd "$(dirname "$0")/.."

npm run --prefix ./frontend watch -- --mode debug &
(cd ./backend/ && air) &
wait
