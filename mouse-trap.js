let activeCircle = null;
let boxElement = null;

function createCircle(event) {
  const circle = document.createElement('div');
  circle.classList.add('circle');
  circle.style.backgroundColor = 'white';
  circle.style.position = 'absolute'; 
  
  const x = event.clientX;
  const y = event.clientY;
  
  circle.style.left = `${x}px`;
  circle.style.top = `${y}px`;
  
  document.body.appendChild(circle);
  activeCircle = circle;

  checkTrap(activeCircle);
}

function moveCircle(event) {
  if (!activeCircle) return;

  const x = event.clientX;
  const y = event.clientY;

  activeCircle.style.left = `${x}px`;
  activeCircle.style.top = `${y}px`;

  checkTrap(activeCircle);
}

function setBox() {
  if (document.querySelector('.box')) {
    boxElement = document.querySelector('.box');
    return;
  }

  boxElement = document.createElement('div');
  boxElement.classList.add('box');
  
  boxElement.style.position = 'fixed';
  boxElement.style.top = '50%';
  boxElement.style.left = '50%';
  boxElement.style.transform = 'translate(-50%, -50%)';
  
  document.body.appendChild(boxElement);
}

function checkTrap(circle) {
  if (!boxElement || !circle) return;

  const circleRect = circle.getBoundingClientRect();
  const boxRect = boxElement.getBoundingClientRect();

  const isStrictlyInside = 
    circleRect.left > boxRect.left &&
    circleRect.right < boxRect.right &&
    circleRect.top > boxRect.top &&
    circleRect.bottom < boxRect.bottom;

  if (isStrictlyInside) {
    circle.style.backgroundColor = 'var(--purple)';
    activeCircle = null; 
  }
}

setBox();
document.addEventListener('click', createCircle);
document.addEventListener('mousemove', moveCircle);