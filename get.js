

const src = { nested: { key: "peekaboo" } };
const path = "nested.key";



function get(src,path){

    if (path===""){
       return undefined
    }

let result=path.split(".");

for (let i=0;i<result.length;i++){

    if (result[i] !==undefined){
        src=src[result[i]]
    }else{
        return undefined;
    }
}


return src;
}

console.log(get(src, path));