function indexOf(arr, value, starts) {
  let start = starts;
  if (start === undefined) {
    start = 0;
  }

  for (let i = start; i < arr.length-1; i++) {
    if (arr[i] === value) {
      return i;
      break;
    }
  }
  return -1;
}

function lastIndexOf(arr, value) {
  for (let j = arr.length; j >= 0; j--) {
    if (arr[j] === value) {
      return j;
      break;
    }
  }
  return -1;
}

function includes(arr, value) {
  for (let i = 0; i < arr.length; i++) {
    if (arr[i] === value) {
      return true;
    }
  }
  return false;
}
