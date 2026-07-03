function pick(obj,keys){

    const keysToPick = Array.isArray(keys) ? keys : [keys];
    return Object.fromEntries(
        Object.entries(obj).filter(([key])=>keysToPick.includes(key))
    )
}

function omit(obj,keys){

    const keysToOmit=Array.isArray(keys) ? keys : [keys]

    return Object.fromEntries(
        Object.entries(obj).filter(([key])=>!keysToOmit.includes(key))
    )
}
console.log(pick({name: "Alice", age: 30, password: "secret"},"age"))
console.log(omit({name: "Alice", age: 30, password: "secret"},"age"))