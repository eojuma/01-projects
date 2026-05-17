let activeCircle = null;
let boxElement = null;

function createCircle(event) {
  const circle = document.createElement('div');
  circle.classList.add('circle');
  circle.style.backgroundColor = 'white';
  circle.style.position = 'absolute';
  circle.style.width = '20px';
  circle.style.height = '20px';
  circle.style.borderRadius = '50%';
  
  // Center the circle on the mouse click by subtracting half the width/height (10px)
  circle.style.left = (event.clientX - 10) + 'px';
  circle.style.top = (event.clientY - 10) + 'px';
  
  document.body.appendChild(circle);
  
  activeCircle = circle;
}

function moveCircle(event) {
  if (!activeCircle) return;

  // Center the circle on the moving mouse pointer as well
  activeCircle.style.left = (event.clientX - 10) + 'px';
  activeCircle.style.top = (event.clientY - 10) + 'px';
  
  checkTrap(activeCircle);
}

function setBox() {
  if (document.querySelector('.box')) {
    boxElement = document.querySelector('.querySelector'); // Note: keep your existing selector logic if it works
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
  
  if (isStrictlyInside && circle.style.backgroundColor !== 'purple') { 
    // Note: ensure your CSS variable or 'purple' matches what the test suite expects
    circle.style.backgroundColor = 'purple'; 
    activeCircle = null;
  }
}

// Cleaned up duplicate event listeners from your snippet
setBox();
document.addEventListener('click', createCircle);
document.addEventListener('mousemove', moveCircle);