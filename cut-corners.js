function round(num){
if (num >= 0) {
    let integerPart = parseInt(num);
    let fractionalPart = num - integerPart;
    if (fractionalPart >= 0.5) {
      return integerPart + 1;
    }
    return integerPart;
  } 
  
  // Handle negative numbers (Note: -1.5 rounds to -1 in JS)
  else {
    let integerPart = parseInt(num);
    let fractionalPart = num - integerPart;
    // For negative numbers, fractionalPart will be negative (e.g., -0.5)
    if (fractionalPart <= -0.51) { // -0.51, -0.6, etc go to -2
      return integerPart - 1;
    }
    return integerPart; // -0.5 goes to -1
  }
}

function ceil(num){
  let integerPart = parseInt(num);
  let fractionalPart = num - integerPart;

  // If there's no decimal, it's already a whole number
  if (fractionalPart === 0) {
    return num;
  }

  // For positive numbers, any decimal means we jump up to the next integer
  if (num > 0) {
    return integerPart + 1;
  }

  // For negative numbers, the integerPart is already "greater" 
  // than the original (e.g., parseInt(-4.1) is -4, and -4 > -4.1)
  return integerPart;
}


function floor(num){
let integerPart = parseInt(num);
  let fractionalPart = num - integerPart;

  if (fractionalPart === 0) {
    return num;
  }

  // For positive numbers, the integerPart is already smaller (4.9 -> 4)
  if (num > 0) {
    return integerPart;
  }

  // For negative numbers, we need to go "down" (e.g., -4.1 -> -5)
  return integerPart - 1;
}

function trunc(num) {
  // If the number is already an integer, return it as is
  // This helps avoid any weirdness with types
  if (num === parseInt(num)) {
    return num;
  }

  // Under your restrictions (no strings, no bitwise, no %), 
  // parseInt is the most direct way to "chop off" the decimal.
  return parseInt(num);
}