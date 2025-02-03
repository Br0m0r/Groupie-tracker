// static/js/filters.js
document.addEventListener('DOMContentLoaded', () => {
    // Get all required elements
    const filterForm = document.getElementById('filterForm');
    const resetButton = document.getElementById('resetFilters');
    
    // Range sliders elements
    const creationYearMin = document.getElementById('creationYearMin');
    const creationYearMax = document.getElementById('creationYearMax');
    const creationYearMinValue = document.getElementById('creationYearMinValue');
    const creationYearMaxValue = document.getElementById('creationYearMaxValue');
    const albumYearMin = document.getElementById('albumYearMin');
    const albumYearMax = document.getElementById('albumYearMax');
    const albumYearMinValue = document.getElementById('albumYearMinValue');
    const albumYearMaxValue = document.getElementById('albumYearMaxValue');
    const memberCheckboxes = document.querySelectorAll('.checkbox-wrapper input[type="checkbox"]');

    // Initialize min/max values
    const currentYear = new Date().getFullYear();
    const minYear = 1958;
    
    // Set initial values
    [creationYearMin, albumYearMin].forEach(slider => {
        slider.min = minYear;
        slider.max = currentYear;
        slider.value = minYear;
    });
    
    [creationYearMax, albumYearMax].forEach(slider => {
        slider.min = minYear;
        slider.max = currentYear;
        slider.value = currentYear;
    });

    // Update display values
    creationYearMinValue.textContent = minYear;
    creationYearMaxValue.textContent = currentYear;
    albumYearMinValue.textContent = minYear;
    albumYearMaxValue.textContent = currentYear;

    // Handle range slider changes
    creationYearMin.addEventListener('input', () => {
        const value = parseInt(creationYearMin.value);
        creationYearMinValue.textContent = value;
        if (value > parseInt(creationYearMax.value)) {
            creationYearMax.value = value;
            creationYearMaxValue.textContent = value;
        }
    });

    creationYearMax.addEventListener('input', () => {
        const value = parseInt(creationYearMax.value);
        creationYearMaxValue.textContent = value;
        if (value < parseInt(creationYearMin.value)) {
            creationYearMin.value = value;
            creationYearMinValue.textContent = value;
        }
    });

    albumYearMin.addEventListener('input', () => {
        const value = parseInt(albumYearMin.value);
        albumYearMinValue.textContent = value;
        if (value > parseInt(albumYearMax.value)) {
            albumYearMax.value = value;
            albumYearMaxValue.textContent = value;
        }
    });

    albumYearMax.addEventListener('input', () => {
        const value = parseInt(albumYearMax.value);
        albumYearMaxValue.textContent = value;
        if (value < parseInt(albumYearMin.value)) {
            albumYearMin.value = value;
            albumYearMinValue.textContent = value;
        }
    });

    // Reset filters
    resetButton.addEventListener('click', (e) => {
        e.preventDefault();
        filterForm.reset();
        
        // Reset slider values
        creationYearMin.value = minYear;
        creationYearMax.value = currentYear;
        albumYearMin.value = minYear;
        albumYearMax.value = currentYear;
        
        // Reset displayed values
        creationYearMinValue.textContent = minYear;
        creationYearMaxValue.textContent = currentYear;
        albumYearMinValue.textContent = minYear;
        albumYearMaxValue.textContent = currentYear;
        
        // Uncheck all member checkboxes
        memberCheckboxes.forEach(checkbox => checkbox.checked = false);
        
        // Apply the reset filters
        applyFilters();
    });

    // Handle form submission
    filterForm.addEventListener('submit', (e) => {
        e.preventDefault();
        applyFilters();
    });

    function applyFilters() {
        const filterData = {
            creationYearMin: parseInt(creationYearMin.value),
            creationYearMax: parseInt(creationYearMax.value),
            albumYearMin: parseInt(albumYearMin.value),
            albumYearMax: parseInt(albumYearMax.value),
            members: Array.from(memberCheckboxes)
                .filter(cb => cb.checked)
                .map(cb => parseInt(cb.value))
        };

        // Show loading state
        const grid = document.querySelector('.artists-grid');
        grid.innerHTML = '<div class="loading">Loading...</div>';

        fetch('/filter', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(filterData)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(artists => {
            grid.innerHTML = ''; // Clear current artists
            
            if (!artists || artists.length === 0) {
                // Show no results message
                const noResults = document.createElement('div');
                noResults.className = 'no-results';
                noResults.textContent = 'No artists found matching your filters';
                grid.appendChild(noResults);
            } else {
                // Create and add artist cards
                artists.forEach(artist => {
                    const card = createArtistCard(artist);
                    grid.appendChild(card);
                });
            }
        })
        .catch(error => {
            console.error('Error:', error);
            grid.innerHTML = '<div class="error-message">Error loading artists. Please try again.</div>';
        });
    }

    function createArtistCard(artist) {
        const card = document.createElement('div');
        card.className = 'artist-card';
        
        const firstAlbumYear = artist.FirstAlbum ? artist.FirstAlbum.match(/\d{4}/)?.[0] : '';
        
        card.innerHTML = `
            <a href="/artist?id=${artist.ID}" class="artist-link">
                <div class="image-container">
                    <img src="${artist.Image}" alt="${artist.Name}" loading="lazy">
                </div>
                <div class="artist-info">
                    <h2>${artist.Name}</h2>
                </div>
            </a>
        `;

        return card;
    }
});