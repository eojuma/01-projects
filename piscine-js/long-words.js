function longWords(arrset) {
    return arrset.every(word => typeof(word) === 'string' && word.length >= 5);
  }

  function oneLongWord(arrset) {
    return arrset.some(word => typeof(word) === 'string' && word.length >= 10);
  }

  function noLongWords(arr) {
    return !arr.some(word => typeof word === 'string' && word.length >= 7);
  }