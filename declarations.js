const escapeStr = "```\\\/\"\'"

const arr  = Object.freeze([4,'2'])

const obj =Object.freeze( {
    str:"joel",
    num:4,
    bool:true,
    undef:undefined,
})

const nested = Object.freeze({
    arr:[4,undefined,'2'],
    obj:{str:"jane",num:9,bool:true,},
})