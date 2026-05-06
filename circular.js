const circular = {};
// setting the object property to match the object itself and the property value
circular["circular"] = circular;
console.log(circular === circular);
