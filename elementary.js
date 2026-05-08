function multiply(a, b) {
  if (a === 0 || b === 0) return 0;
  let isNegative = (a < 0 && b > 0) || (a > 0 && b < 0);

  if (a < 0) a = -a;
  if (b < 0) b = -b;

  let result = 0;

  for (let i = 0; i < b; i++) {
    result += a;
  }

  if (isNegative) {
    return -result;
  }
  return result;
}

function divide(a, b) {
  let count = 0;
  if (b === 0) {
    return "undefined";
  }

  let isNegative = a < 0 && b < 0;

  if (a < 0) {
    a = -a;
  }

  if (b < 0) {
    b = -b;
  }

  while (a >= b) {
    a -= b;
    count++;
  }

  if (isNegative) {
    return -count;
  }
  return count;
}

function modulo(a, b) {
  if (b === 0) {
    return "NaN";
  }

  let isNegative = a < 0;

  if (a < 0) {
    a = -a;
  }
  if (b < 0) {
    b = -b;
  }

  while (a >= b) {
    a -= b;
  }

  if (isNegative) {
    return -a;
  }
  return a;
}

console.log(multiply(-8, 8));
console.log(divide(12, 2));
console.log(modulo(0, 0));
