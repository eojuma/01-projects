function slice(input, starts, ends) {
  let start = starts;
  let end = ends;
  if (start === undefined) {
    start = 0;
  }

  if (end === undefined) {
    end = input.length;
  }
  if (start < 0) {
    start = input.length + start;
  }
  if (end < 0) {
    end = input.length + end;
  }
  let result = [];

  for (let i = start; i < end; i++) {
    result.push(input[i]);
  }

  if (typeof input === "string") {
    return result.join("");
  }
  return result;
}
