function filerShortStateName(states){
    return states.filter(state=>state.length<7);
}

function filterStartVowel(words) {
    return words.filter(str => /^[aeiou]/i.test(str));
  }