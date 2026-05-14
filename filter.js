function filter(arrset, callback) {
    const result = [];
    for (let i = 0; i < arrset.length; i++) {
      if (callback(arrset[i], i, arrset)) {
        result.push(arrset[i]);
      }
    }
    return result;
  }

  function reject(arrset, callback) {
    const result = [];
    for (let i = 0; i < arrset.length; i++) {
      if (!callback(arrset[i], i, arrset)) {
        result.push(arrset[i]);
      }
    }
    return result;
  }