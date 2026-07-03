function map(arr, callback) {
    const result = [];
    for (let i = 0; i < arr.length; i++) {
      result.push(callback(arr[i], i, arr));
    }
    return result;
  }

  function flatMap(arrset, callback) {
    const result = [];
    
    for (let i = 0; i < arrset.length; i++) {
      const mappervalue = callback(arrset[i], i, arrset);
      
      if (Array.isArray(mappervalue)) {
        for (let j = 0; j < mappervalue.length; j++) {
          result.push(mappervalue[j]);
        }
      } else {
        result.push(mappervalue);
      }
    }
    
    return result;
  }