// 1. is.num: Must exclude NaN based on instructions 
// (If test #1 failed previously, it was likely due to missing !Number.isNaN)
is.num = (n) => typeof n === 'number' && !Number.isNaN(n);

// 2. is.nan: True only for NaN
is.nan = (n) => Number.isNaN(n);

// 3. is.str: True for strings
is.str = (s) => typeof s === 'string';

// 4. is.bool: True for booleans
is.bool = (b) => typeof b === 'boolean';

// 5. is.undef: True for undefined
is.undef = (u) => typeof u === 'undefined';

// 6. is.def: True if NOT undefined
is.def = (d) => typeof d !== 'undefined';

// 7. is.arr: Uses the Array-specific check
is.arr = (a) => Array.isArray(a);

// 8. is.obj: The fix for Test #7
// We check for 'object', exclude null, and exclude arrays.
is.obj = (o) => typeof o === 'object' && o !== null && !Array.isArray(o);

// 9. is.fun: True for functions
is.fun = (f) => typeof f === 'function';

// 10. is.truthy: Simple boolean coercion
is.truthy = (t) => !!t;

// 11. is.falsy: Simple boolean negation
is.falsy = (f) => !f;