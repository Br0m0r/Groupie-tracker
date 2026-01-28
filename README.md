# Groupie Tracker

Groupie Tracker is a Go web app that visualizes artists, their concert locations, and dates using the Groupie Trackers API.

## Features
- Artist list with detail pages
- Search with live suggestions (artist, member, location, dates)
- Filters by creation date, first album year, member count, and locations
- Concert locations map (Leaflet + OpenStreetMap) with geocoding via Nominatim

## Requirements
- Go 1.22+
- Internet access (API data, map tiles, geocoding)

## Run
```bash
go run main.go
```
Open `http://localhost:8080`

## Routes
- `GET /` home page
- `GET /artist?id={id}` artist details
- `GET /search?q={query}` search page / suggestions
- `GET /filter` filtered results
- `GET /api/coordinates?id={id}` artist map coordinates (JSON)

## Project Structure
```
main.go
handlers/
models/
store/
utils/
templates/
static/
```
