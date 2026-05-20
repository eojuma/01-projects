// A strict check to ensure we only recursively merge plain data objects
const isPlainObject = (item) => {
  return (
    item !== null &&
    typeof item === 'object' &&
    Object.prototype.toString.call(item) === '[object Object]'
  );
};

function replica(target, ...sources) {
  // If target isn't a plain object, initialize it as one
  if (!isPlainObject(target)) {
    target = {};
  }

  for (const source of sources) {
    if (!isPlainObject(source)) continue;

    for (const key in source) {
      if (Object.prototype.hasOwnProperty.call(source, key)) {
        const sourceValue = source[key];

        if (isPlainObject(sourceValue)) {
          if (!isPlainObject(target[key])) {
            target[key] = {};
          }
          replica(target[key], sourceValue);
        } else if (Array.isArray(sourceValue)) {
          target[key] = sourceValue.map(item => 
            isPlainObject(item) || Array.isArray(item) ? replica(Array.isArray(item) ? [] : {}, item) : item
          );
        } else if (sourceValue instanceof RegExp) {
          target[key] = new RegExp(sourceValue.source, sourceValue.flags);
        } else if (sourceValue instanceof Date) {
          target[key] = new Date(sourceValue.getTime());
        } else {
          target[key] = sourceValue;
        }
      }
    }
  }

  return target;
}