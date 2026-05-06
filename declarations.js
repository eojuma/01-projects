const escapeStr = "`\\\\/\"'";

const arr = Object.freeze([4, "2"]);

const obj = Object.freeze({
  str: "joel",
  num: 4,
  bool: true,
  undef: undefined,
});

const nestedInnerObj = Object.freeze({
  str: "jane",
  num: 9,
  bool: true,
});

const nestedInnerArr = Object.freeze([4, undefined, "2"]);

const nested = Object.freeze({
  arr: nestedInnerArr,
  obj: nestedInnerObj,
});
