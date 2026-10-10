# Groupie Tracker

**Collaborative School Project**
This project was carried out as a group during our computer science studies at Ynov Campus.

## Project Overview
Groupie Tracker is a desktop application we developed in Go (Golang) using the Fyne graphical toolkit. It interacts with an external REST API to display information about various musical artists and bands (history, members, concert dates and locations).

The goal of this project was to familiarize ourselves with data manipulation, user interface design in Go, and asynchronous API request handling (including geocoding and map visualization).

## Features
- **Artist Directory**: Display artists in a grid format with their names and images.
- **Advanced Search**: We implemented a search bar to filter by artist name, member name, creation date, first album year, or concert location.
- **Filter System**: The list can be refined by creation year, number of members, etc.
- **Detailed View**: Displays complete information for a selected artist (members, albums, concerts).
- **Geolocation & Map**:
  - The application converts location names into geographical coordinates via the Nominatim API.
  - We integrated an interactive map (OpenStreetMap) to visualize concert locations.

## Tech Stack
- **Language**: Go (Golang)
- **GUI Framework**: Fyne v2
- **Data Format**: JSON
- **External APIs**:
  - Artist data (API provided in the assignment)
  - Geocoding: Nominatim (OpenStreetMap)
  - Map tiles: OpenStreetMap

## Project Structure
```
projet_groupie-Tracker/
├── main.go            # Application entry point
├── appli/             # Core logic and UI components
│   ├── api.go         # HTTP request handling and data structures
│   └── page.go        # Fyne layout, events, and map rendering
├── go.mod / go.sum    # Go module dependencies
└── README.md          # Documentation
```

## Prerequisites
To run this project, you will need:
1. **Go** (version 1.25 or compatible).
2. **A C compiler**: Fyne requires a C compiler (GCC) for CGO bindings related to graphical rendering.
   - *Windows*: TDM-GCC or MinGW-w64.
   - *macOS*: Xcode Command Line Tools.
   - *Linux*: GCC.

## Installation and Launch
1. Clone this repository.
2. Open a terminal in the root folder.
3. Download the dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

### 👥 Contributors
- Ynov Campus Students Group (Go Project).
