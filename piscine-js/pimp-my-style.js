import { styles } from './pimp-my-style.data.js';

let index = 0;
let isRemoving = false;

export function pimp() {
  const button = document.querySelector('.button');

  if (!isRemoving) {
   
    button.classList.add(styles[index]);
    index++;
    if (index === styles.length) {
      isRemoving = true;
      button.classList.add('unpimp');
    }
  } else {
    index--;
    button.classList.remove(styles[index]);

    if (index === 0) {
      isRemoving = false;
      button.classList.remove('unpimp');
    }
  }
}