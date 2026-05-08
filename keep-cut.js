

function cutFirst(s){
    return s.slice(2)
}

function cutLast(s){
    return s.slice(0,-2)
}

function cutFirstLast(s){
    return s.slice(2,-2)
}
function keepFirst(s){
    return s.slice(0,2)
}
function keepLast(s){
    return s.slice(-2)
}

function keepFirstLast(s) {
  if (s.length <= 3) {
    return s;
  }
  return s.slice(0, 2) + s.slice(-2);
}

