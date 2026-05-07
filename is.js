const is = {
  num: function numbered(n) {
    if (typeof n === "number") {
      return true;
    }
    return false;
  },
  nan: function nan(a) {
    if (typeof a === "NaN") {
      return true;
    }
    return false;
  },
  str: function str(s) {
    if (typeof s === "string") {
      return true;
    }
    return false;
  },
  bool: function bool(b) {
    if (typeof b === "boolean") {
      return true;
    }
    return false;
  },
  undef: function undef(u) {
    if (typeof u === "undefined") {
      return true;
    }
    return false;
  },
  arr: function def(d) {
    if (typeof d === "undefined") {
      return false;
    }
    return true;
  },
  obj: function object(o) {
    if (typeof o === "null") {
      return true;
    }
    return false;
  },
  fun: function fun(f) {
    if (typeof f === "function") {
      return true;
    }
    return false;
  },
  truthy: function truth(t) {
    if (typeof t === "truthy") {
      return true;
    }
    return false;
  },
  falsy: function fals(f) {
    if (typeof f === "falsy") {
      return true;
    }
    return false;
  },
};
