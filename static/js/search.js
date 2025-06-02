document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.querySelector('.search-input');
    const suggestionsList = document.querySelector('.suggestions-list');
    const searchContainer = document.querySelector('.search-suggestions');

    // Minimal clean approach
    searchInput.addEventListener('input', (e) => {
        const query = e.target.value.trim();
        
        if (query === '') {
            hide();
            return;
        }

        // Fetch suggestions
        fetch(`/search?q=${encodeURIComponent(query)}`, {
            headers: { 'X-Requested-With': 'XMLHttpRequest' }
        })
        .then(response => response.json())
        .then(results => {
            if (results && results.length > 0) {
                show(results);
            } else {
                hide();
            }
        })
        .catch(() => hide());
    });

    function show(results) {
        const html = results.map(result => `
            <div class="suggestion-item">
                <div class="suggestion-text">${result.Text}</div>
                <div class="suggestion-type">${result.Type}</div>
            </div>
        `).join('');
        
        suggestionsList.innerHTML = html;
        searchContainer.style.display = 'block';
    }

    function hide() {
        searchContainer.style.display = 'none';
        suggestionsList.innerHTML = '';
    }

    // Hide on outside click
    document.addEventListener('click', (e) => {
        if (!searchContainer.contains(e.target) && !searchInput.contains(e.target)) {
            hide();
        }
    });

    // Hide on Escape
    searchInput.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') hide();
    });
});