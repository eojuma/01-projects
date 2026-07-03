function fusion(...obj) {
  return obj.reduce((acc, current) => {
    const keysAvailable = [
      ...new Set([...Object.keys(acc), ...Object.keys(current)]),
    ];
    const result = {};

    for (const key of keysAvailable) {
      const v1 = acc[key];
      const v2 = current[key];

      if (v1 === undefined) {
        result[key] = v2;
        continue;
      }

      if (v2 === undefined) {
        result[key] = v1;
        continue;
      }

      const t1 = typeof v1;
      const t2 = typeof v2;

      if (Array.isArray(v1) && Array.isArray(v2)){
        result[key]=[...v1,...v2];
      }else if(Array.isArray(v1) || Array.isArray(v2)){
        result[key]=v2;
      }else if (t1 === "number" && t2 === "number") {
        result[key] = v1 + v2;
      } else if (t1 === "string" && t2 === "string") {
        result[key] = v1 + " " + v2;
      } else if (t1 === "object" && t2 === "object") {
        result[key] = fusion(v1, v2);
      }else{
        result[key]=v2;
      }
    }
    return result;
  },{});
}

console.log(fusion({ arr: [1, "2"] }, { arr: [2] }));
console.log(fusion({ a: { b: [3, 2], c: { d: 8 } } }, { a: { b: [0, 3, 1], c: { d: 3 } } }))
console.log(fusion({ a: "hello", b: [] }, { a: 4 }))