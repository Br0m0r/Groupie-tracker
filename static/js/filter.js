document.addEventListener('DOMContentLoaded', () => {
    const filterToggle = document.querySelector('.filter-toggle');
    const filterContent = document.querySelector('.filter-content');
    const applyFiltersBtn = document.getElementById('applyFilters');
    const resetFiltersBtn = document.getElementById('resetFilters');
    const locationSearch = document.getElementById('locationSearch');
    
    // Initial state - hidden
    filterContent.style.display = 'none';
    
    // Show/Hide filters
    filterToggle.addEventListener('click', () => {
        const isHidden = filterContent.style.display === 'none';
        filterContent.style.display = isHidden ? 'block' : 'none';
        filterToggle.textContent = isHidden ? 'Hide Filters' : 'Show Filters';
        
        if (isHidden) {
            setTimeout(() => filterContent.classList.add('active'), 10);
        } else {
            filterContent.classList.remove('active');
        }
    });

    // Range input value display updates
    const creationStartInput = document.getElementById('creationDateStart');
    const creationEndInput = document.getElementById('creationDateEnd');
    const albumStartInput = document.getElementById('firstAlbumStart');
    const albumEndInput = document.getElementById('firstAlbumEnd');

    // Update displayed values when range inputs change
    creationStartInput.addEventListener('input', () => {
        document.getElementById('creationStartValue').textContent = creationStartInput.value;
    });

    creationEndInput.addEventListener('input', () => {
        document.getElementById('creationEndValue').textContent = creationEndInput.value;
    });

    albumStartInput.addEventListener('input', () => {
        document.getElementById('albumStartValue').textContent = albumStartInput.value;
    });

    albumEndInput.addEventListener('input', () => {
        document.getElementById('albumEndValue').textContent = albumEndInput.value;
    });

    // Location search functionality
    locationSearch.addEventListener('input', (e) => {
        const searchValue = e.target.value.toLowerCase();
        const locationCheckboxes = document.querySelectorAll('#locationFilter label');
        
        locationCheckboxes.forEach(label => {
            const labelText = label.textContent.toLowerCase();
            label.style.display = labelText.includes(searchValue) ? '' : 'none';
        });
    });

    // Apply filters
    applyFiltersBtn.addEventListener('click', () => {
        const filters = {
            creationStart: parseInt(creationStartInput.value),
            creationEnd: parseInt(creationEndInput.value),
            albumStart: parseInt(albumStartInput.value),
            albumEnd: parseInt(albumEndInput.value),
            members: Array.from(document.querySelectorAll('input[name="members"]:checked'))
                        .map(cb => parseInt(cb.value)),
            locations: Array.from(document.querySelectorAll('input[name="locations"]:checked'))
                        .map(cb => cb.value)
        };

        // Send filter request
        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/', true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');

        xhr.onload = function() {
            if (xhr.status === 200) {
                try {
                    const response = JSON.parse(xhr.responseText);
                    updateArtistsGrid(response.artists);
                    document.getElementById('resultCount').textContent = response.artists.length;
                } catch (error) {
                    console.error('Error parsing response:', error);
                }
            } else {
                console.error('Filter request failed:', xhr.status);
            }
        };

        xhr.onerror = function() {
            console.error('Filter request failed');
        };

        xhr.send(JSON.stringify(filters));
    });

    // Reset filters
    resetFiltersBtn.addEventListener('click', () => {
        // Reset range inputs
        creationStartInput.value = '1950';
        creationEndInput.value = '2024';
        albumStartInput.value = '1950';
        albumEndInput.value = '2024';
        
        // Update displayed values
        document.getElementById('creationStartValue').textContent = '1950';
        document.getElementById('creationEndValue').textContent = '2024';
        document.getElementById('albumStartValue').textContent = '1950';
        document.getElementById('albumEndValue').textContent = '2024';
        
        // Clear location search
        locationSearch.value = '';
        
        // Show all location checkboxes
        document.querySelectorAll('#locationFilter label').forEach(label => {
            label.style.display = '';
        });
        
        // Uncheck all checkboxes
        document.querySelectorAll('input[type="checkbox"]').forEach(cb => {
            cb.checked = false;
        });
        
        // Apply filters to reset view
        document.getElementById('applyFilters').click();
    });
});

// Helper function to update the artists grid
function updateArtistsGrid(artists) {
    const artistsGrid = document.querySelector('.artists-grid');
    
    console.log('Updating grid with artists:', artists); // Debug log

    if (!artists || artists.length === 0) {
        artistsGrid.innerHTML = '<div class="no-results">No artists found matching your filters</div>';
        return;
    }

    const html = artists.map(artist => {
        console.log('Processing artist:', artist); // Debug individual artist
        return `
            <div class="artist-card">
                <a href="/artist?id=${artist.ID || artist.id}" class="artist-link">
                    <div class="image-container">
                        <img src="${artist.Image || artist.image}" alt="${artist.Name || artist.name}">
                    </div>
                    <div class="artist-info">
                        <h2>${artist.Name || artist.name}</h2>
                    </div>
                </a>
            </div>
        `;
    }).join('');

    console.log('Generated HTML:', html); // Debug generated HTML
    artistsGrid.innerHTML = html;
}

function applyFilters() {
    const filters = {
        creationStart: parseInt(document.getElementById('creationDateStart').value),
        creationEnd: parseInt(document.getElementById('creationDateEnd').value),
        albumStart: parseInt(document.getElementById('firstAlbumStart').value),
        albumEnd: parseInt(document.getElementById('firstAlbumEnd').value),
        members: Array.from(document.querySelectorAll('input[name="members"]:checked'))
                    .map(cb => parseInt(cb.value)),
        locations: Array.from(document.querySelectorAll('input[name="locations"]:checked'))
                    .map(cb => cb.value)
    };

    // Send filter request
    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');

    xhr.onload = function() {
        if (xhr.status === 200) {
            try {
                const response = JSON.parse(xhr.responseText);
                if (response && response.artists) {
                    updateArtistsGrid(response.artists);
                    document.getElementById('resultCount').textContent = response.artists.length;
                }
            } catch (error) {
                console.error('Error parsing response:', error);
            }
        }
    };

    xhr.send(JSON.stringify(filters));
}
