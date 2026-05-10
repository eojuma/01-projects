var smalls = -Infinity;
var biggie = Infinity;

//finding the biggest number in an array
const nums =[39,7,10]
let biggest =-Infinity

for (let n of nums){
    if (n>biggest){
        biggest=n;
    }
}

//finding the smallest number in array

let smallest=Infinity

numbers=[89,9,8]

for (let t of numbers){
    if (t<smallest){
        smallest=t;
    }
}
console.log(biggest);
console.log(smallest);