function letterSpaceNumber(s) {
    let result = [];

    for (let i = 0; i < s.length - 2; i++) {
        const char = s[i];
        const space = s[i+1];
        const digit = s[i+2];
        const nextChar = s[i+3]; 

        if (
            space === ' ' && 
            (digit >= '0' && digit <= '9') && 
            !(nextChar >= '0' && nextChar <= '9') 
        ) {
            result.push(char + space + digit);
        }
    }
    return result;
}