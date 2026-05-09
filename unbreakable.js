function split(str, separator) {
  let result = "";
  let final = [];

  for (let i = 0; i < str.length; i++) {
    if (str[i] !== separator) {
      result += str[i];
    } else {
      final.push(result);
      result = "";
    }
  }
  
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
