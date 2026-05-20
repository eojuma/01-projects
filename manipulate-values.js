function filterValues(obj,callback){

    return Object.fromEntries(
Object.entries(obj).filter(([KeyboardEvent,value])=>callback(value,KeyboardEvent,obj))
    );
}

function mapValues(obj,callback){
    return Object.fromEntries(
        Object.entries(obj).map(([Key,value])=>[Key,callback(value,Key,obj)])
    );
}


function reduceValues(obj,callback,initialValue){
     return Object.entries(obj).reduce((acc, [key, value]) => {
    return callback(acc, value, key, obj);
  }, initialValue);
}