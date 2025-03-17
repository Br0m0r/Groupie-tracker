Overview

Groupie Tracker integrates data from a third-party music band API to present a comprehensive catalog of artists, including their images, detailed biographies, concert locations, and performance dates. The application is built with Go for the backend, using standard libraries, and vanilla HTML/CSS/JavaScript for the frontend. Key functionalities include a dynamic search that provides real-time suggestions, an interactive map feature built with OpenStreetMap and Leaflet.js, and a filtering system that lets users narrow down results by criteria such as creation dates, album release years, member counts, and concert locations.
Features
Artist Information

    Catalog View: Displays a grid of artist cards with images and names.
    Detailed Artist Pages: 
        Each artist page shows:
           Basic artist details (name, image, creation date, first album)
           Band members list
           Concert locations and dates
           Relation mapping between locations and performance dates
    Error Handling: Comprehensive error pages for invalid requests or server issues.

Map Feature

    Interactive Map: Visualizes concert locations using OpenStreetMap and Leaflet.js.
    Marker Interaction: Clickable markers display location details.
    Geocoding: Converts addresses to geographic coordinates using the Nominatim API with rate limiting.
    Responsive Design: The map adjusts to different screen sizes for an optimal viewing experience.

Dynamic Search

    Real-Time Suggestions: Provides search suggestions as the user types, with categories for artists, band members, locations, creation dates, and first album details.
    AJAX-Enabled: Utilizes AJAX for dynamic query processing and suggestion updates without full page reloads.
    Result Types: Each result displays the main text, type indicator, description, and direct link to the respective artist page.

Filtering

    Customizable Filters: Users can filter artists based on:
        Number of band members
        Creation date range
        First album year range
        Concert locations
    Responsive UI: The filter panel is designed to be mobile-friendly and provides real-time result counts.
    Default Values: The application automatically sets sensible default filter parameters based on the dataset.

Tech Stack

    Backend: Go (using only standard packages)
    Frontend: HTML, CSS, JavaScript
    Map: OpenStreetMap Tiles, Leaflet.js (v1.9.3)
    Geocoding: Nominatim API
    Data Consumption: REST API endpoints provided by the third-party music band API




Setup & Running
Prerequisites

    Go 1.22 or higher

Installation:

    

    Use the following URL to clone the repo:

    git clone https://github.com/Br0m0r/Groupie-tracker.git


Start the application using:

    go run .

    Access the Application

    Open your browser and navigate to http://localhost:8080.

Development Notes

    Error Handling: The application includes robust error handling for invalid inputs, failed API calls, and template parsing errors.
    Responsive Design: The UI adapts to various screen sizes, ensuring a smooth experience on both desktop and mobile devices.
    Rate Limiting: Geocoding requests to the Nominatim API are rate-limited to adhere to usage policies.
    Data Caching: Coordinates are cached in the background to reduce repeated API calls and improve performance.
    Dynamic Search: The real-time search functionality is built with AJAX, ensuring seamless user interaction without page reloads.

