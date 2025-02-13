# Groupie Tracker Filters

A web application that displays and filters information about music bands and artists, including their locations, concert dates, and other details.

## Overview

This project is a web application that allows users to explore and filter information about various music artists and bands. It features an intuitive interface with multiple filtering options and a sophisticated search functionality.

## Features

### Search Functionality
- Real-time search suggestions
- Search across multiple categories:
  - Artist/Band names
  - Band members
  - Concert locations
  - Creation dates
  - First album dates

### Advanced Filtering System
1. **Member Count Filter**
   - Checkbox-based filtering
   - Options from 1 to 8+ members
   - Multiple selection supported

2. **Creation Date Range**
   - Slider-based range selection
   - Range from 1950 to 2024
   - Visual feedback with real-time updates

3. **First Album Date Range**
   - Similar slider-based interface
   - Range from 1950 to 2024
   - Interactive sliders with value display

4. **Location Filter**
   - Hierarchical location structure
   - Special handling for US locations
   - State-City relationship mapping

### Special Location Handling for USA

The application implements a sophisticated geography system for US locations, recognizing the hierarchical relationship between states and cities. This is managed through the `geography.go` utility:

```go
var stateCityMap = map[string][]string{
    "Washington":     {"Seattle"},
    "California":     {"Los Angeles", "Anaheim", "Oakland", "Del Mar", "San Francisco"},
    "Texas":          {"Dallas", "Houston"},
    // ... other states and cities
}
```

This system allows for:
- Automatic state recognition from city names
- Grouping cities by their respective states
- Hierarchical filtering of locations
- Intelligent location search suggestions

## Technical Implementation

### Backend
- Written in Go (1.22.2)
- Uses only standard Go packages
- Implements concurrent data fetching
- RESTful API architecture

### Frontend
- Pure HTML/CSS/JavaScript
- Responsive design
- Real-time interaction
- No external frontend frameworks

### Key Components
1. **Data Store**
   - Concurrent data fetching
   - Thread-safe operations
   - Efficient data caching

2. **Error Handling**
   - Comprehensive error types
   - User-friendly error pages
   - Detailed error logging

3. **Utils Package**
   - Location formatting
   - Date formatting
   - Geography management

## Project Structure

```
groupie/
├── handlers/
│   ├── errors.go
│   ├── handlers.go
│   ├── search.go
│   └── filter.go
├── models/
│   ├── artistModel.go
│   ├── filterModel.go
│   └── searchModel.go
├── store/
│   └── store.go
├── utils/
│   ├── utils.go
│   └── geography.go
├── static/
│   ├── css/
│   └── js/
├── templates/
│   ├── artist.html
│   ├── error.html
│   ├── index.html
│   └── search.html
└── main.go
```

## Running the Project

1. Ensure Go 1.22.2 or later is installed
2. Clone the repository
3. Navigate to the project directory
4. Run the application:
   ```bash
   go run main.go
   ```
5. Access the application at `http://localhost:8080`

## API Integration

The application integrates with the Groupie Tracker API (`https://groupietrackers.herokuapp.com/api`) to fetch:
- Artist information
- Concert locations
- Performance dates
- Relation data

## Performance Optimizations

1. **Concurrent Data Fetching**
   - Uses Go routines for parallel API requests
   - Implements efficient error handling
   - Manages concurrent access to shared resources

2. **Caching**
   - In-memory data storage
   - Thread-safe operations
   - Efficient data retrieval

3. **Frontend Optimization**
   - Lazy loading of images
   - Efficient DOM updates
   - Debounced search functionality

## Error Handling

The application implements a comprehensive error handling system:
- Custom error types for different scenarios
- User-friendly error pages
- Detailed error logging
- Graceful degradation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is part of the 01.edu system.