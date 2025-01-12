// static/js/main.js
document.addEventListener('DOMContentLoaded', () => {
    // Create background container
    const backgroundContainer = document.createElement('div');
    backgroundContainer.className = 'background-animation';
    document.body.prepend(backgroundContainer);

    // Music-related SVG objects
    const musicObjects = {
        note: `<svg viewBox="0 0 24 24"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>`,
        doubleNote: `<svg viewBox="0 0 24 24"><path d="M19 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h-4V3h2zm-9 0v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h-4V3h2z"/></svg>`,
        headphones: `<svg viewBox="0 0 24 24"><path d="M12 1c-4.97 0-9 4.03-9 9v7c0 1.66 1.34 3 3 3h3v-8H5v-2c0-3.87 3.13-7 7-7s7 3.13 7 7v2h-4v8h3c1.66 0 3-1.34 3-3v-7c0-4.97-4.03-9-9-9z"/></svg>`,
        guitar: `<svg viewBox="0 0 24 24"><path d="M19.59 3H22v2h-1.59l-3.5 3.5c-.19.19-.44.29-.71.29-.55 0-1-.45-1-1 0-.27.1-.52.29-.71L19.59 3zM9 3v2H7v2H5v2H3v2h2v2h2v2h2v2h2v2h2v-2h2v-2h2v-2h2v-2h-2v-2h-2V9h-2V7h-2V5h-2V3H9z"/></svg>`,
        microphone: `<svg viewBox="0 0 24 24"><path d="M12 14c1.66 0 3-1.34 3-3V5c0-1.66-1.34-3-3-3S9 3.34 9 5v6c0 1.66 1.34 3 3 3zm-1-9c0-.55.45-1 1-1s1 .45 1 1v6c0 .55-.45 1-1 1s-1-.45-1-1V5zm6 6c0 2.76-2.24 5-5 5s-5-2.24-5-5H5c0 3.53 2.61 6.43 6 6.92V21h2v-3.08c3.39-.49 6-3.39 6-6.92h-2z"/></svg>`,
        vinyl: `<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="3"/></svg>`,
        speaker: `<svg viewBox="0 0 24 24"><path d="M17 2H7c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h10c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 18H7V4h10v16z"/></svg>`
    };

    // Create initial set of floating objects
    const TOTAL_OBJECTS = 40;

    function createFloatingObject() {
        const object = document.createElement('div');
        const objectTypes = Object.keys(musicObjects);
        const type = objectTypes[Math.floor(Math.random() * objectTypes.length)];
        const size = Math.random() * 30 + 20;
        
        object.className = 'floating-object';
        object.style.width = `${size}px`;
        object.style.height = `${size}px`;
        object.innerHTML = musicObjects[type];
        
        // Random positions across the entire viewport
        object.style.left = `${Math.random() * 100}vw`;
        object.style.top = `${Math.random() * 100}vh`;
        
        // Set random animation duration
        const duration = Math.random() * 8 + 7; // 7-15s duration
        object.style.animationDuration = `${duration}s`;
        
        return object;
    }

    // Create initial set of objects
    for (let i = 0; i < TOTAL_OBJECTS; i++) {
        backgroundContainer.appendChild(createFloatingObject());
    }
});