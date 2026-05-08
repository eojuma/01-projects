function getInteger(n) {
  if (n === 0) return 0;
  let isNeg = n < 0;
  let num = isNeg ? -n : n;
  
  let result = 0;
  let step = 1000000000; // Start with huge steps
  
  while (step >= 1) {
    while (result + step <= num) {
      result += step;
    }
    step /= 10; // Reduce step size (1000 -> 100 -> 10 -> 1)
  }
  
  return isNeg ? -result : result;
}

// Now the functions are instant!
function trunc(num) {
  return getInteger(num);
}

function floor(num) {
  let int = getInteger(num);
  return (num < 0 && num !== int) ? int - 1 : int;
}

function ceil(num) {
  let int = getInteger(num);
  return (num > 0 && num !== int) ? int + 1 : int;
}

function round(num) {
  let int = getInteger(num);
  let frac = num - int;
  if (num >= 0) {
    return frac >= 0.5 ? int + 1 : int;
  } else {
    // Special check for the -1.5 -> -1 JS behavior
    return frac <= -0.51 ? int - 1 : int;
  }
}