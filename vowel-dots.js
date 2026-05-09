function vowelDots(s) {
  let result = "";
  for (let i = 0; i < s.length; i++) {
    if (isVowel(s[i])) {
      result += s[i] + ".";
    } else {
      result += s[i];
    }
  }

  return result;
}

function isVowel(s) {
  return (
    s === "a" ||
    s === "e" ||
    s === "i" ||
    s === "o" ||
    s === "u" ||
    s === "A" ||
    s === "E" ||
    s === "I" ||
    s === "O" ||
    s === "U"
  );
}
