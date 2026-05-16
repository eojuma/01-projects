import { colors } from './fifty-shades-of-cold.data.js';

export function generateClasses() {
  const styleTag = document.createElement('style');
  let cssRules = '';

  colors.forEach(color => {
    cssRules += `.${color} { background: ${color}; }\n`;
  });

  styleTag.textContent = cssRules;
  document.head.appendChild(styleTag);
}

export function generateColdShades() {
  const coldKeywords = ['aqua', 'blue', 'turquoise', 'green', 'cyan', 'navy', 'purple'];

  colors.forEach(color => {
    const isCold = coldKeywords.some(keyword => color.includes(keyword));

    if (isCold) {
      const div = document.createElement('div');
      div.className = color;
      div.textContent = color;
      
      div.addEventListener('click', () => choseShade(color));
      
      document.body.appendChild(div);
    }
  });
}

export function choseShade(shade) {
  const allDivs = document.querySelectorAll('div');
  
  allDivs.forEach(div => {
    div.className = shade;
  });
}