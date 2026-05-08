function getInteger(n) {
  if (n === 0) return 0;
  let i = 0;
  if (n > 0) {
    while (i + 1 <= n) {
      i++;
    }
    return i;
  } else {
    while (i - 1 >= n) {
      i--;
    }
    return i;
  }
}

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
    // JS Math.round logic: -1.5 -> -1
    return frac <= -0.51 ? int - 1 : int;
  }
}