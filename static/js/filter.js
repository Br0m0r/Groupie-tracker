    // filter.js
    function extractYear(dateString) {
        // Handle date format "DD-MM-YYYY"
        if (dateString.includes('-')) {
            return parseInt(dateString.split('-')[2]);
        }
        // If it's already a year
        return parseInt(dateString);
    }
    document.addEventListener('DOMContentLoaded', () => {
        // Filter Panel Toggle
        const filterToggle = document.getElementById('filterToggle');
        const filterPanel = document.getElementById('filterPanel');

        filterToggle.addEventListener('click', () => {
            filterPanel.classList.toggle('active');
        });

        // Close filter panel when clicking outside
        document.addEventListener('click', (e) => {
            if (!filterPanel.contains(e.target) && !filterToggle.contains(e.target)) {
                filterPanel.classList.remove('active');
            }
        });

        // Year Range Sliders
        function setupRangeSlider(minId, maxId, minValueId, maxValueId) {
            const minSlider = document.getElementById(minId);
            const maxSlider = document.getElementById(maxId);
            const minValue = document.getElementById(minValueId);
            const maxValue = document.getElementById(maxValueId);

            function updateRange() {
                if (parseInt(minSlider.value) > parseInt(maxSlider.value)) {
                    const tmp = maxSlider.value;
                    maxSlider.value = minSlider.value;
                    minSlider.value = tmp;
                }
                minValue.textContent = minSlider.value;
                maxValue.textContent = maxSlider.value;
                applyFilters();
            }

            minSlider.addEventListener('input', updateRange);
            maxSlider.addEventListener('input', updateRange);
        }

        setupRangeSlider('creationYearMin', 'creationYearMax', 'creationYearMinValue', 'creationYearMaxValue');
        setupRangeSlider('albumYearMin', 'albumYearMax', 'albumYearMinValue', 'albumYearMaxValue');

        // Locations Toggle
        const toggleLocations = document.getElementById('toggleLocations');
        const locationsContainer = document.getElementById('locationsContainer');

        toggleLocations.addEventListener('click', () => {
            locationsContainer.classList.toggle('hidden');
            const svg = toggleLocations.querySelector('svg');
            svg.style.transform = locationsContainer.classList.contains('hidden') 
                ? 'rotate(0deg)' 
                : 'rotate(180deg)';
        });

        // Member checkboxes
        const memberCheckboxes = document.querySelectorAll('.member-checkbox');
        memberCheckboxes.forEach(checkbox => {
            checkbox.addEventListener('change', applyFilters);
        });

        // Function to collect all unique locations from artist cards
        function collectLocations() {
            const locations = new Set();
            document.querySelectorAll('.artist-card').forEach(card => {
                const artistLocations = card.dataset.locations;
                if (artistLocations) {
                    artistLocations.split(',').forEach(location => {
                        locations.add(location.trim());
                    });
                }
            });
            return Array.from(locations).sort();
        }
        
        // Populate locations
        function populateLocations() {
            const locationsList = document.querySelector('.locations-list');
            const locations = collectLocations();
            
            locationsList.innerHTML = locations.map(location => `
                <label class="location-checkbox">
                    <input type="checkbox" value="${location}">
                    <span>${location}</span>
                </label>
            `).join('');

            // Add change event listeners to new checkboxes
            locationsList.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
                checkbox.addEventListener('change', applyFilters);
            });
        }

        // Initialize locations
        populateLocations();

        // Filter Application
        function applyFilters() {
            const creationYearMin = parseInt(document.getElementById('creationYearMin').value);
            const creationYearMax = parseInt(document.getElementById('creationYearMax').value);
            const albumYearMin = parseInt(document.getElementById('albumYearMin').value);
            const albumYearMax = parseInt(document.getElementById('albumYearMax').value);
            
            const selectedMembers = Array.from(document.querySelectorAll('.member-checkbox:checked'))
                .map(checkbox => parseInt(checkbox.value));
            
            const selectedLocations = Array.from(document.querySelectorAll('.location-checkbox input:checked'))
                .map(checkbox => checkbox.value);

            const artistCards = document.querySelectorAll('.artist-card');
            
            artistCards.forEach(card => {
                const creationYear = parseInt(card.dataset.creationYear);
                const firstAlbum = extractYear(card.dataset.firstAlbum);
                const memberCount = parseInt(card.dataset.members);
                const artistLocations = card.dataset.locations ? card.dataset.locations.split(',').map(l => l.trim()) : [];
            
                const creationYearMatch = creationYear >= creationYearMin && creationYear <= creationYearMax;
                const albumYearMatch = firstAlbum >= albumYearMin && firstAlbum <= albumYearMax;
                const memberMatch = selectedMembers.length === 0 || selectedMembers.includes(memberCount);
                const locationMatch = selectedLocations.length === 0 || 
                    selectedLocations.some(location => artistLocations.includes(location));
            
                if (creationYearMatch && albumYearMatch && memberMatch && locationMatch) {
                    card.style.display = '';
                    card.style.animation = 'fadeInUp 0.5s ease forwards';
                } else {
                    card.style.display = 'none';
                }
            });
        }

        // Add animation keyframes if not already in your CSS
        if (!document.querySelector('#filterAnimations')) {
            const style = document.createElement('style');
            style.id = 'filterAnimations';
            style.textContent = `
                @keyframes fadeInUp {
                    from {
                        opacity: 0;
                        transform: translateY(20px);
                    }
                    to {
                        opacity: 1;
                        transform: translateY(0);
                    }
                }
            `;
            document.head.appendChild(style);
        }
    });

    