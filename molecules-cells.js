function RNA(s) {
  let result = "";

  for (let i = 0; i < s.length; i++) {
    if (s[i] === "G") {
      result += "C";
    } else if (s[i] === "C") {
      result += "G";
    } else if (s[i] === "T") {
      result += "A";
    } else if (s[i] === "A") {
      result += "U";
    }
  }
  return result;
}

function DNA(s) {
  let result = "";
  for (let i = 0; i < s.length; i++) {
    if (s[i] === "C") {
      result += "G";
    } else if (s[i] === "G") {
      result += "C";
    } else if (s[i] === "A") {
      result += "T";
    } else if (s[i] === "U") {
      result += "A";
    }
  }
  return result;
}
