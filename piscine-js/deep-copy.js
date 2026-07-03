function deepCopy(val) {
  if (val === null || typeof val !== 'object') {
    return val;
  }

  if (val instanceof Date) return new Date(val.getTime());
  if (val instanceof RegExp) return new RegExp(val.source, val.flags);

  const copy = Array.isArray(val) ? [] : {};

  for (const key in val) {
    if (Object.prototype.hasOwnProperty.call(val, key)) {
      copy[key] = deepCopy(val[key]);
    }
  }

  return copy;
}