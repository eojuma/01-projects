function flat(arr) {
  const result = [];
  
  function flatten(item) {
    if (Array.isArray(item)) {
      for (let element of item) {
        flatten(element);
      }
    } else {
      result.push(item);
    }
  }
  
  flatten(arr);
  return result;
}

console.log(flat([1, [2, [3, [4]]]])); // [1, 2, 3, 4]
console.log(flat([1, [2], [[3]], 4])); // [1, 2, 3, 4]
console.log(flat([])); // []
console.log(flat([1, 2, 3])); // [1, 2, 3]
console.log(flat([[[]]])); // []