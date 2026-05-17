let activeCircle = null;
let boxElement = null;

function createCircle(event) {
  const circle = document.createElement('div');
  circle.classList.add('circle');
  circle.style.backgroundColor = 'white';
  circle.style.position = 'fixed';
  circle.style.left = event.clientX + 'px';
  circle.style.top = event.clientY + 'px';
  circle.style.width = '20px';
  circle.style.height = '20px';
  circle.style.borderRadius = '50%';

  document.body.appendChild(circle);
  activeCircle = circle;
}

function moveCircle(event) {
  if (!activeCircle) return;

  activeCircle.style.left = event.clientX + 'px';
  activeCircle.style.top = event.clientY + 'px';

  checkTrap(activeCircle);
}

function setBox() {
  if (document.querySelector('.box')) {
    boxElement = document.querySelector('.box');
    return;
  }

  boxElement = document.createElement('div');
  boxElement.classList.add('box');
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