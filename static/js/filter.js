// static/js/filter.js
document.addEventListener('DOMContentLoaded', () => {
    // Get the clear button and filter form
    const clearButton = document.querySelector('.clear-filters');
    const filterForm = document.getElementById('filter-form');

    if (clearButton && filterForm) {
        clearButton.addEventListener('click', (e) => {
            e.preventDefault(); // Prevent the default reset behavior

            // Clear all checkboxes
            filterForm.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
                checkbox.checked = false;
            });

            // Clear all number inputs and set to default values
            filterForm.querySelectorAll('input[type="number"]').forEach(input => {
                if (input.name.includes('creation')) {
                    input.value = input.name.includes('start') ? '1950' : '2024';
                } else if (input.name.includes('album')) {
                    input.value = input.name.includes('start') ? '1950' : '2024';
                }
            });

            // Clear search input if exists
            const searchInput = filterForm.querySelector('.location-search-input');
            if (searchInput) {
                searchInput.value = '';
            }

            // Submit the form to refresh with cleared filters
            filterForm.submit();
        });
    }
});