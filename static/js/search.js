document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.querySelector('.search-input');
    const suggestionsList = document.querySelector('.suggestions-list');
    const searchContainer = document.querySelector('.search-suggestions');

    // Simple search suggestions - just visual, no clicking
    searchInput.addEventListener('input', (e) => {
        const query = e.target.value.trim();
        
        // Hide suggestions if query is empty
        if (query === '') {
            hideResults();
            return;
        }

        // Show loading state (optional)
        showResults();
        suggestionsList.innerHTML = '<div class="loading">Searching...</div>';

        // Make AJAX request for suggestions
        fetch(`/search?q=${encodeURIComponent(query)}`, {
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        })
        .then(response => response.json())
        .then(results => {
            displayResults(results);
        })
        .catch(error => {
            console.error('Search error:', error);
            suggestionsList.innerHTML = '<div class="error">Search failed</div>';
        });
    });

    // Display the search results
    function displayResults(results) {
        if (!results || results.length === 0) {
            suggestionsList.innerHTML = '<div class="no-results">No suggestions found</div>';
            return;
        }

        // Create simple HTML for each result
        const html = results.map(result => `
            <div class="suggestion-item">
                <div class="suggestion-text">${result.Text || 'Unknown'}</div>
                <div class="suggestion-type">${result.Type || 'unknown'}</div>
            </div>
        `).join('');

        suggestionsList.innerHTML = html;
    }

    // Show the suggestions container
    function showResults() {
        searchContainer.style.display = 'block';
    }

    // Hide the suggestions container
    function hideResults() {
        searchContainer.style.display = 'none';
        suggestionsList.innerHTML = '';
    }

    // Hide suggestions when clicking outside
    document.addEventListener('click', (e) => {
        if (!searchContainer.contains(e.target) && !searchInput.contains(e.target)) {
            hideResults();
        }
    });

    // Hide suggestions when pressing Escape
    searchInput.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            hideResults();
        }
    });
});