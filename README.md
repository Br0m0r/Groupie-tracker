# Groupie Tracker

A web application that visualizes artist information and concert locations using data from an external API.

## Features

### ✅ Core Functionality
- **Artist Visualization**: Display artist cards with information (name, image, creation year, first album, members)
- **Concert Locations**: Interactive map showing where artists have performed
- **Real-time Search**: Live search suggestions for artists, members, locations, and dates
- **Advanced Filtering**: Filter by creation date, album date, member count, and locations
- **Responsive Design**: Modern UI that works on all devices

### ✅ Search System
- Search by artist/band name, members, locations, first album date, creation date
- Case-insensitive search with live suggestions
- Type identification in suggestions (artist, member, location, etc.)

### ✅ Filter System
- **Range filters**: Creation date (1958-2015), First album date (1963-2020)
- **Checkbox filters**: Member count (1-8 members), Concert locations
- Real-time filtering without page reload

### ✅ Geolocation Map
- Interactive world map using OpenStreetMap and Leaflet.js
- Address-to-coordinates conversion via Nominatim API
- Clickable markers showing concert locations
- Responsive map design

## Quick Start

### Prerequisites
- Go 1.19 or later
- Internet connection (for API data and map tiles)

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd groupie-tracker

# Run the application
go run main.go
```

### Usage
1. Open your browser and navigate to `http://localhost:8080`
2. Browse artists on the homepage
3. Use the search bar for quick artist lookup
4. Apply filters to narrow down results
5. Click on any artist to view detailed information and concert map

## Project Structure
```
├── main.go              # Server entry point
├── handlers/            # HTTP request handlers
│   ├── handlers.go      # Main handlers
│   ├── search.go        # Search functionality
│   ├── filter.go        # Filter functionality
│   └── coordinates.go   # Map coordinates API
├── models/              # Data structures
├── store/               # Data storage and API client
├── utils/               # Helper functions
├── templates/           # HTML templates
└── static/              # CSS, JavaScript, and assets
    ├── css/            # Stylesheets
    └── js/             # Client-side scripts
```

## API Endpoints
- `GET /` - Homepage with artist listings
- `GET /artist?id={id}` - Artist details page with map
- `GET /search?q={query}` - Search results
- `GET /filter` - Filtered artist results
- `GET /api/coordinates?id={id}` - Artist location coordinates (JSON)

## Technologies Used
- **Backend**: Go (standard library only)
- **Frontend**: HTML5, CSS3, Vanilla JavaScript
- **Map**: Leaflet.js + OpenStreetMap
- **Geocoding**: Nominatim API
- **Data Source**: Groupie Trackers API

## Error Handling
- Graceful error pages for 404, 400, and 500 errors
- Rate limiting for external API calls
- Fallback behavior for failed geocoding requests