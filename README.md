# DMX Viewer

This project is a web application to view DMX data, built with Go and React.

## Getting Started

### Prerequisites

- Go
- Node.js and npm
- air (Go hot-reloading tool)

### Installation and Running

1. **Install frontend dependencies:**
   ```bash
   npm install --prefix frontend
   ```

2. **Build the React app:**
   ```bash
   npm run build --prefix frontend
   ```

3. **Install air:**
   ```bash
   go install github.com/air-verse/air@latest
   ```

4. **Run the Go server with hot-reloading:**
   ```bash
   cd backend
   air
   ```

5. **Open your browser and navigate to [http://localhost:8080](http://localhost:8080)**

## Directory Structure

This project's main directory structure is as follows:

```
dmx_viewer/
├── backend/              # Backend service implemented in Go
├── docs/                 # Project documentation
├── frontend/             # Frontend application implemented in React
├── scripts/              # 開発用のユーティリティスクリプト
```
