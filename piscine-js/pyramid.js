function pyramid(s, n) {
  if (n <= 0 || s === ""){

      return "";
  }

  let result = "";
  let sLen = s.length;

  for (let i = 1; i <= n; i++) {
    let spaces = " ".repeat((n - i) * sLen);

    let bricks = s.repeat(2 * i - 1);

    result += spaces + bricks;

    if (i < n) {
      result += "\n";
    }
  }

  return result;
}

console.log(pyramid("*", 5));