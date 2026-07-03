function ionOut(s){
    const regex = /\b(\w+t)(?=ion\b)/g;

    const matches = s.match(regex);

    return matches || [];
}