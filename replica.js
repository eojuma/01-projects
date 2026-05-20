const isObject = (item) => {
  return item && typeof item === 'object' && !Array.isArray(item);
};

function replica(target, ...sources) {
  if (!isObject(target)) {
    target = {};
  }

  for (const source of sources) {
    if (!isObject(source)) continue;

    for (const key in source) {
      if (Object.prototype.hasOwnProperty.call(source, key)) {
        const sourceValue = source[key];

        if (isObject(sourceValue)) {
          if (!isObject(target[key])) {
            target[key] = {};
          }
          replica(target[key], sourceValue);
        } else if (Array.isArray(sourceValue)) {
          target[key] = sourceValue.map(item => 
            isObject(item) || Array.isArray(item) ? replica(Array.isArray(item) ? [] : {}, item) : item
          );
        } else {
          target[key] = sourceValue;
        }
      }
    }
  }

  return target;
}