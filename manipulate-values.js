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
    const values = Object.values(obj);
    if (initialValue === undefined) {
    return values.reduce(callback);
  }
  
  return values.reduce(callback, initialValue);
}