#!/bin/zsh -e

cd "$(dirname "$0")/.."

npm run --prefix ./frontend watch &
(cd ./backend/ && air) &
wait
