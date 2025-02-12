document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.querySelector('.search-input');
    const suggestionsList = document.querySelector('.suggestions-list');
    const searchContainer = document.querySelector('.search-suggestions');

   
    searchInput.addEventListener('input', (e) => {
        const query = e.target.value.trim();
        
        if (query === '') {
            suggestionsList.innerHTML = '';
            searchContainer.style.display = 'none';
            return;
        }

        fetch(`/search?q=${encodeURIComponent(query)}`, {
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        })
        .then(response => response.json())
        .then(results => {
            if (results.length === 0) {
                suggestionsList.innerHTML = '';
                searchContainer.style.display = 'none';
                return;
            }

            searchContainer.style.display = 'block';
            // Render all results without limiting to 5
            suggestionsList.innerHTML = results.map(suggestion => `
                <div class="suggestion-item">
                    <span class="suggestion-text">${suggestion.text}</span>
                    <span class="suggestion-type">${suggestion.type}</span>
                </div>
            `).join('');
        })
        .catch(error => {
            console.error('Failed to fetch suggestions:', error);
            searchContainer.style.display = 'none';
        });
    });

    // Close suggestions when clicking outside
    document.addEventListener('click', (e) => {
        if (!searchContainer.contains(e.target) && !searchInput.contains(e.target)) {
            searchContainer.style.display = 'none';
        }
    });
});