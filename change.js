const sourceObject = {
  num: 42,
  bool: true,
  str: "some text",
  log: console.log,
};


function get(key){
    return sourceObject[key];
}


function set(keys,values){
    return sourceObject[keys]=values;
}

// console.log(get("bool"))
// console.log(set("str","john"))