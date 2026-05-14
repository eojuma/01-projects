function filerShortStateName(states){
    return states.filter(state=>state.length<7);
}

function filterStartVowel(words) {
    return words.filter(str => /^[aeiou]/i.test(str));
  }

  function filter5Vowels(words) {
    return words.filter(str => {
      const vowels = str.match(/[aeiou]/gi) || [];
      return vowels.length >= 5;
    });
  }

  function filter1DistinctVowel(words) {
    return words.filter(str => {
      const vowels = str.match(/[aeiou]/gi) || [];
      if (vowels.length === 0) return false;
      
      const distinctVowels = new Set(vowels.map(v => v.toLowerCase()));
      return distinctVowels.size === 1;
    });
  }