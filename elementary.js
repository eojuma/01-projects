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

  // Case: Both negative (Result should be positive)
  if (a < 0 && b < 0) {
    a = -a;
    b = -b;
    while (a >= b) {
      a -= b;
      count++;
    }
    return count; // Correct: - / - = +
  }

  // Case: a is negative, b is positive (Result should be negative)
  if (a < 0 && b > 0) {
    a = -a;
    while (a >= b) {
      a -= b;
      count++;
    }
    return -count; // FIX: - / + = -
  }

  // Case: b is negative, a is positive (Result should be negative)
  if (b < 0 && a > 0) {
    b = -b;
    while (a >= b) {
      a -= b;
      count++;
    }
    return -count; // FIX: + / - = -
  }

  // Case: Both positive (Result should be positive)
  if (a > 0 && b > 0) {
    while (a >= b) {
      a -= b;
      count++;
    }
    return count; // FIX: Added return here to avoid the "0" trap at the end
  }

  // Handle zero
  if (a === 0) return 0;

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
