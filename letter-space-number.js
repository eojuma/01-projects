function letterSpaceNumber(s) {
    let result = [];

    for (let i = 0; i < s.length - 2; i++) {
        if (s[i+1] === ' ' && s[i+2] >= '0' && s[i+2] <= '9') {
            result.push(s[i] + s[i+1] + s[i+2])+"\n";
            break; 
        }
    }
    return result;
}

console.log(letterSpaceNumber("example 1, example 20"));