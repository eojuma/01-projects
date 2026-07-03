function arrToSet(val) {
  return new Set(val);
}

function arrToStr(val) {
  return val.join('');
}

function setToArr(val) {
  return [...val];
}

function setToStr(val) {
  return [...val].join('');
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
  return Object.values(val);
}

function objToMap(val) {
  return new Map(Object.entries(val));
}

function arrToObj(val) {
  return { ...val };
}

function strToObj(val) {
  return { ...val };
}

function superTypeOf(val) {
  if (val === null) return 'null';
  if (val === undefined) return 'undefined';
  return Object.prototype.toString.call(val).slice(8, -1);
}