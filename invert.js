function invert(obj){


const inverted = Object.fromEntries(
    Object.entries(obj).map(([Key,val])=>[val,Key])
);
return inverted
}

console.log(invert({ a: "x", b: "y", c: "z" }))