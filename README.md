# DMX Viewer

## Overview

DMX Viewer provides a web page to visualize ArtNet signals in real time.

![dmx_viewer](https://github.com/user-attachments/assets/da3a52e7-0712-4b6e-8bbc-7249f7bb9c4e)

## Features

- **Display DMX signals in real time**  
Displays received ArtNet signals in tables and graphs.
- **Show received universes per node**  
Displays the state of universes received by each node.
- **Check changes with history graphs**  
Displays changes in DMX values over time.

## Architecture

- **Frontend**
	- Framework: React (TypeScript)
	- Backend: Go
	- Communication: WebSocket
