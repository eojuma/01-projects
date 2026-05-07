// const is = {
  const num= function numbered(n) {
    if (typeof n === "number") {
      return true;
    }
    return false;
  }
  const nan= function nan(a) {
    if (typeof a === "NaN") {
      return true;
    }
    return false;
  }
  const str= function str(s) {
    if (typeof s === "string") {
      return true;
    }
    return false;
  }
 const bool= function bool(b) {
    if (typeof b === "boolean") {
      return true;
    }
    return false;
  }
  const undef= function undef(u) {
    if (typeof u === "undefined") {
      return true;
    }
    return false;
  }
 const arr= function def(d) {
    if (typeof d === "undefined") {
      return false;
    }
    return true;
  }
const  obj= function object(o) {
    if (typeof o === "null") {
      return true;
    }
    return false;
  }
 const fun= function fun(f) {
    if (typeof f === "function") {
      return true;
    }
    return false;
  }
 const truthy= function truth(t) {
    if (typeof t === "truthy") {
      return true;
    }
    return false;
  }
 const falsy= function fals(f) {
    if (typeof f === "falsy") {
      return true;
    }
    return false;
  }
// };
