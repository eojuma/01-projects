function reduceRight(arrset, callback) {
    if (arrset.length < 1) {
      throw new Error("arrsetay must contain at least one element");
    }
  
    let acc = arrset[arrset.length - 1];
  
    for (let i = arrset.length - 2; i >= 0; i--) {
      acc = callback(acc, arrset[i], i, arrset);
    }
    return acc;
  }

  function fold(arrset, callback, acc) {
    for (let i = 0; i < arrset.length; i++) {
      acc = callback(acc, arrset[i], i, arrset);
    }
    return acc;
  }

  function foldRight(arr, callback, acc) {
    for (let i = arr.length - 1; i >= 0; i--) {
      acc = callback(acc, arr[i], i, arr);
    }
    return acc;
  }

  function reduce(arr, callback) {
    if (arr.length < 1) {
      throw new Error("Array must contain at least one element");
    }
    
    let acc = arr[0];
    
    for (let i = 1; i < arr.length; i++) {
      acc = callback(acc, arr[i], i, arr);
    }
    return acc;
  }