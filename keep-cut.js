

function cutFirst(s){
    return s.slice(2)
}

function cutLast(s){
    return s.slice(0,s.length-2)
}

function cutFirstLast(s){
    return s.slice(1,s.length-1)
}
function keepFirst(s){
    return s.slice(0,1)
}
function keepLast(s){
    return s.slice(s.length-2,s.length-1)
}

function keepFirstLast(s){
    return s.slice(2,-2)
}

