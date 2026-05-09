const vowels = /[aeiou]/gi;
function vowelDots(s) {
  
    return s.replace(vowels,'$&.');
}

