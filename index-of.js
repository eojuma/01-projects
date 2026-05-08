function indexOf(arr, value, starts) {
  let start = 0;
  if (start === "undefined") {
    start = starts;
  }

  for (let i = start; i < arr.length; i++) {
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
