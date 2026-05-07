// is.is = {
is.num = function numbered(n) {
  if (typeof n === "number") {
    return true;
  }
  return false;
};

is.nan = function (a) {
  return Number.isNaN(a);
};

is.str = function str(s) {
  if (typeof s === "string") {
    return true;
  }
  return false;
};
is.bool = function bool(b) {
  if (typeof b === "boolean") {
    return true;
  }
  return false;
};
is.undef = function undef(u) {
  if (typeof u === "undefined") {
    return true;
  }
  return false;
};
is.arr = function def(d) {
  if (typeof d === "undefined") {
    return false;
  }
  return true;
};
is.obj = function object(o) {
  if (typeof o === "null") {
    return true;
  }
  return false;
};
is.fun = function fun(f) {
  if (typeof f === "function") {
    return true;
  }
  return false;
};
is.truthy = function truth(t) {
  if (typeof t === "truthy") {
    return true;
  }
  return false;
};
is.falsy = function fals(f) {
  if (typeof f === "falsy") {
    return true;
  }
  return false;
};
// };
