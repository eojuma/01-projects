function compose() {
    document.addEventListener('keydown', (event) => {
        const key = event.key;

        if (/^[a-z]$/.test(key)) {
            const note = document.createElement('div');
            note.classList.add('note');
            note.textContent = key;

            note.style.backgroundColor = generateColorFromKey(key);

            document.body.appendChild(note);
        }
        
        else if (key === 'Backspace') {
            const notes = document.querySelectorAll('.note');
            if (notes.length > 0) {
                const lastNote = notes[notes.length - 1];
                lastNote.remove();
            }
        }
        
        else if (key === 'Escape') {
            const notes = document.querySelectorAll('.note');
            notes.forEach(note => note.remove());
        }
    });
}

function generateColorFromKey(key) {
    const charCode = key.charCodeAt(0);
    
    const hue = ((charCode - 97) * (360 / 26)) % 360;
    
    return `hsl(${hue}, 75%, 60%)`;
}

if (typeof window !== 'undefined') {
    compose();
}