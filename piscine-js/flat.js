function flat(arr, depth = 1) {
  if (depth === 0) return arr;

  const result = [];

  for (let item of arr) {
    if (Array.isArray(item)) {
      result.push(...flat(item, depth - 1));
    } else {
      result.push(item);
    }
  }

  return result;
}