export function generateLetters() {
    const container = document.getElementById('container') || document.body;
    
    const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    const totalLetters = 120;

    for (let i = 1; i <= totalLetters; i++) {
        const letterDiv = document.createElement('div');

        const randomIndex = Math.floor(Math.random() * alphabet.length);
        letterDiv.textContent = alphabet[randomIndex];

        const fontSize = 11 + (i - 1);
        letterDiv.style.fontSize = `${fontSize}px`;

        if (i <= 40) {
            letterDiv.style.fontWeight = "300"; 
        } else if (i <= 80) {
            letterDiv.style.fontWeight = "400"; 
        } else {
            letterDiv.style.fontWeight = "600"; 
        }

        container.appendChild(letterDiv);
    }
}