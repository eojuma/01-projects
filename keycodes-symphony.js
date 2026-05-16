export function compose() {
    document.addEventListener('keydown', (event) => {
        const key = event.key;

        if (key >= 'a' && key <= 'z' && key.length === 1) {
            const note = document.createElement('div');
            note.classList.add('note');
            note.textContent = key;
            note.style.backgroundColor = generateColorFromKey(key);
            document.body.appendChild(note);
        }
        
        else if (key === 'Backspace') {
            const notes = document.querySelectorAll('.note');
            if (notes.length > 0) {
                notes[notes.length - 1].remove();
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