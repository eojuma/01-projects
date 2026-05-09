function reverse(input){
    let result = [];


    for (let i=input.length-1;i>=0;i--){
     result.push(input[i]);
    }

    if (typeof(input==="string")){
        
            return result.join("");
        
    }
    return result;
}