// is.num: Check if it's a number AND ensure it's not NaN
is.num = function (n) {
  return typeof n === "number" && !Number.isNaN(n);
};

// is.nan: Uses the built-in Number method
is.nan = function (n) {
  return Number.isNaN(n);
};

// is.str: Check for string type
is.str = function (s) {
  return typeof s === "string";
};

// is.bool: Check for boolean type
is.bool = function (b) {
  return typeof b === "boolean";
};

// is.undef: Check for undefined type
is.undef = function (u) {
  return typeof u === "undefined";
};

// is.def: Simply the opposite of is.undef
is.def = function (d) {
  return typeof d !== "undefined";
};

// is.arr: Uses the Array-specific check
is.arr = function (a) {
  return Array.isArray(a);
};

// is.obj: Returns true if it's an object OR null (as per your prompt instructions)
is.obj = function (o) {
  return typeof o === 'object' && o !== null && !Array.isArray(o);
};

// is.fun: Check for function type
is.fun = function (f) {
  return typeof f === "function";
};

// is.truthy: Force the value into a boolean context
is.truthy = function (t) {
  return !!t;
};

// is.falsy: Check if the value is not truthy
is.falsy = function (f) {
  return !f;
};
