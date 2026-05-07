function arrToSet(val) {
  return new Set(val);
}

function arrToStr(val) {
  return val.join("");
}

function setToArr(val) {
  return [...val];
}

function setToStr(val) {
  return [...val].join("");
}

function strToArr(val) {
  return [...val];
}

function strToSet(val) {
  return new Set(val);
}

function mapToObj(val) {
  return Object.fromEntries(val);
}

function objToArr(val) {
  return Object.entries(val);
}

function objToMap(val) {
  return new Map(Object.entries(val));
}

function arrToObj(val) {
  return Object.fromEntries(val); 
}

function strToObj(val) {
  return JSON.parse(val); 
}

function superTypeOf(val) {
  if (val === null) return "null";
  const typeTag = Object.prototype.toString.call(val);
  return typeTag.slice(8, -1); // Returns 'Array', 'Set', etc.
}