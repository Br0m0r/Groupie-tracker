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

    // Handle range slider changes
    creationYearMin.addEventListener('input', () => {
        creationYearMinValue.textContent = creationYearMin.value;
        if (parseInt(creationYearMin.value) > parseInt(creationYearMax.value)) {
            creationYearMax.value = creationYearMin.value;
            creationYearMaxValue.textContent = creationYearMin.value;
        }
    });

    creationYearMax.addEventListener('input', () => {
        creationYearMaxValue.textContent = creationYearMax.value;
        if (parseInt(creationYearMax.value) < parseInt(creationYearMin.value)) {
            creationYearMin.value = creationYearMax.value;
            creationYearMinValue.textContent = creationYearMax.value;
        }
    });

    albumYearMin.addEventListener('input', () => {
        albumYearMinValue.textContent = albumYearMin.value;
        if (parseInt(albumYearMin.value) > parseInt(albumYearMax.value)) {
            albumYearMax.value = albumYearMin.value;
            albumYearMaxValue.textContent = albumYearMin.value;
        }
    });

    albumYearMax.addEventListener('input', () => {
        albumYearMaxValue.textContent = albumYearMax.value;
        if (parseInt(albumYearMax.value) < parseInt(albumYearMin.value)) {
            albumYearMin.value = albumYearMax.value;
            albumYearMinValue.textContent = albumYearMax.value;
        }
    });

    // Reset filters
    resetButton.addEventListener('click', (e) => {
        e.preventDefault();
        filterForm.reset();
        
        // Reset values
        creationYearMin.value = "1958";
        creationYearMax.value = "2024";
        albumYearMin.value = "1958";
        albumYearMax.value = "2024";
        
        // Reset displayed values
        creationYearMinValue.textContent = "1958";
        creationYearMaxValue.textContent = "2024";
        albumYearMinValue.textContent = "1958";
        albumYearMaxValue.textContent = "2024";
        
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

        fetch('/filter', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(filterData)
        })
        .then(response => response.json())
        .then(artists => {
            const grid = document.querySelector('.artists-grid');
            grid.innerHTML = ''; // Clear current artists
            
            if (artists.length === 0) {
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
            const grid = document.querySelector('.artists-grid');
            grid.innerHTML = '<div class="error-message">Error loading artists</div>';
        });
    }

    function createArtistCard(artist) {
        const card = document.createElement('div');
        card.className = 'artist-card';
        card.dataset.creationDate = artist.CreationDate;
        card.dataset.firstAlbum = artist.FirstAlbum;
        card.dataset.members = artist.Members.length;

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

        // Add animation
        setTimeout(() => {
            card.style.animation = 'fadeInUp 0.6s ease forwards';
        }, 0);

        return card;
    }
});