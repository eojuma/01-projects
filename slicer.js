function slice(input, starts, ends) {
  let start = starts;
  let end = ends;
  if (start === undefined) {
    start = 0;
  }

  if (end === undefined) {
    end = input.length;
  }

  let result = [];

  for (let i = start; i < end; i++) {
    result.push(input[i]);
  }

  if (typeof input === "string") {
    result.join("");
  }
  return result;
}
