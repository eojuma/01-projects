function chunk(arr, size) {
  let result = [];
  if (size<1){
    return [];
  }
  for (let i = 0; i < arr.length; i += size) {
    let start = i;
    let end = size + i;

    if (end > arr.length) {
      end = arr.length;
    }
    result.push(arr.slice(start,end));
  }

  return result;
}
