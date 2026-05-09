
function split(str, separator) {
  let final = [];
  let result = "";
  
  // If the separator is empty, JS split('') returns an array of every character
  if (separator === "") {
    for (let i = 0; i < str.length; i++) {
      final.push(str[i]);
    }
    return final;
  }

  for (let i = 0; i < str.length; i++) {
    // We check if the characters starting at 'i' match the separator
    let match = true;
    for (let j = 0; j < separator.length; j++) {
      if (str[i + j] !== separator[j]) {
        match = false;
        break;
      }
    }

    if (match) {
      // If we found the whole separator, push the word we've built
      final.push(result);
      result = "";
      // Move 'i' forward by the length of the separator 
      // (-1 because the loop's i++ will add one more)
      i += separator.length - 1;
    } else {
      // No match? Just add the single character to our current word
      result += str[i];
    }
  }

  // Always push the remaining part after the loop
  final.push(result);
  return final;
}
function join(arr, separator) {
  let result = "";

  for (let i = 0; i < arr.length; i++) {
    // 1. Add the current element
    result += arr[i];

    // 2. Add the separator ONLY if it's not the last element
    // We check (arr.length - 1) to see if we are at the end
    if (i < arr.length - 1) {
      result += separator;
    }
  }

  return result;
}
