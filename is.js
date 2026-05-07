// is.num: The test explicitly expects NaN to return true.
is.num = (n) => typeof n === 'number';

// is.nan: Returns true if the value is NaN.
is.nan = (n) => Number.isNaN(n);

// is.str: Returns true if the value is a string.
is.str = (s) => typeof s === 'string';

// is.bool: Returns true if the value is a boolean.
is.bool = (b) => typeof b === 'boolean';

// is.undef: Returns true if the value is undefined.
is.undef = (u) => typeof u === 'undefined';

// is.def: Returns true if the value is defined.
is.def = (d) => typeof d !== 'undefined';

// is.arr: Returns true if the value is an array.
is.arr = (a) => Array.isArray(a);

// is.obj: The test expects false for null and false for arrays.
is.obj = (o) => typeof o === 'object' && o !== null && !Array.isArray(o);

// is.fun: Returns true if the value is a function.
is.fun = (f) => typeof f === 'function';

// is.truthy: Returns true if the value is truthy.
is.truthy = (t) => !!t;

// is.falsy: Returns true if the value is falsy.
is.falsy = (f) => !f;