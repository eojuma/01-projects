function filerKeys(obj,callback){
const result = {};

for (const key in obj){
if (Object.prototype.hasOwnProperty.call(obj, key)) {
      if (callback(key, obj[key], obj)) {
        result[key] = obj[key];
      }
    }
  }
  return result;
}

function mapKeys(obj, callback) {
  const result = {};
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      const newKey = callback(key, obj[key], obj);
      result[newKey] = obj[key];
    }
  }
  return result;
}

function reduceKeys(obj, callback, initialValue) {
  const keys = Object.keys(obj);
  
  if (initialValue === undefined) {
    return keys.reduce(callback);
  }
  
  return keys.reduce(callback, initialValue);
}