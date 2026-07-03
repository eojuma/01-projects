function triangle(s,n){

    if (n<0  || s===""){
        return "";
    }

let brick = "";

for (let l=1;l<=n;l++){

    brick+=s.repeat(l);

    if (l<n){
        brick+="\n"
    }
}

    return brick;
}

console.log(triangle("n", 5));