


function get(src,path){

let result=path.split(".");

for (let i=0;i<result.length;i++){

    if (src===null || src===undefined){
        return undefined;
    }else{
        src=src[result[i]]
    }
}


return src;
}

console.log(get(src, path));