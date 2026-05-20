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

