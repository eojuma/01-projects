function arrToSet(val) {
  return new Set(val);
}

function arrToStr(val) {
  return val.join(" ");
}

function setToArr(val) {
  return [...val];
}

function setToStr(val) {
  return val.toString();
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

function objToMap(arr) {
  return new Map(Object.entries(val));
}

function arrToObj(arr) {
  return Object.fromEntries(val);
}

function strToObj(arr) {
  return JSON.parse(val);
}

function superTypeOf(arr) {
  if (val === null) {
    return "null";
  }

  const typeTag = Object.prototype.toString.call(val);
  return typeTag.slice(8, -1);
}
