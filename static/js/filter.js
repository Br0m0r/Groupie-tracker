// static/js/filters.js
document.addEventListener('DOMContentLoaded', () => {
    const filterForm = document.getElementById('filterForm');
    const resetButton = document.getElementById('resetFilters');
    const artistCards = document.querySelectorAll('.artist-card');
    
    // Range sliders
    const creationYearRange = document.getElementById('creationYearRange');
    const creationYearValue = document.getElementById('creationYearValue');
    const albumYearRange = document.getElementById('albumYearRange');
    const albumYearValue = document.getElementById('albumYearValue');
    
    // Member checkboxes
    const memberCheckboxes = document.querySelectorAll('.checkbox-wrapper input[type="checkbox"]');

    // Initialize ranges
    function updateYearLabel(range, valueSpan) {
        valueSpan.textContent = range.value;
    }

    creationYearRange.addEventListener('input', () => 
        updateYearLabel(creationYearRange, creationYearValue));
    
    albumYearRange.addEventListener('input', () => 
        updateYearLabel(albumYearRange, albumYearValue));

    // Reset filters
    resetButton.addEventListener('click', (e) => {
        e.preventDefault();
        filterForm.reset();
        creationYearValue.textContent = creationYearRange.value;
        albumYearValue.textContent = albumYearRange.value;
        applyFilters();
    });

    // Apply filters on form submit
    filterForm.addEventListener('submit', (e) => {
        e.preventDefault();
        applyFilters();
    });

    function applyFilters() {
        const filters = {
            creationYear: parseInt(creationYearRange.value),
            albumYear: parseInt(albumYearRange.value),
            members: Array.from(memberCheckboxes)
                .filter(cb => cb.checked)
                .map(cb => parseInt(cb.value))
        };

        artistCards.forEach(card => {
            const creationDate = parseInt(card.dataset.creationDate);
            const firstAlbum = card.dataset.firstAlbum;
            const memberCount = parseInt(card.dataset.members);
            
            let visible = true;

            // Check creation date
            if (filters.creationYear && creationDate > filters.creationYear) {
                visible = false;
            }

            // Check first album (if it's a valid year)
            const albumYear = parseInt(firstAlbum);
            if (!isNaN(albumYear) && filters.albumYear && albumYear > filters.albumYear) {
                visible = false;
            }

            // Check number of members
            if (filters.members.length > 0) {
                const matchesMemberCount = filters.members.some(count => {
                    if (count === 8) {
                        return memberCount >= 8;
                    }
                    return memberCount === count;
                });
                if (!matchesMemberCount) {
                    visible = false;
                }
            }

            // Apply visibility with animation
            if (visible) {
                card.style.display = '';
                card.style.animation = 'fadeInUp 0.6s ease forwards';
            } else {
                card.style.display = 'none';
            }
        });
    }

    // Initial load
    applyFilters();
});