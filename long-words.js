function longWords(arrset) {
    return arrset.every(word => typeof(word) === 'string' && word.length >= 5);
  }