function split(s) {
  let result = "";
  let final = [];

  for (let a = 0; a < s.length; a++) {
    if (s[a] != " ") {
      result += s[a];
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
