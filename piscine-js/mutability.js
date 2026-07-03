// const person = {
//     name: "Rick",
//     age: 77,
//     country: "US",
//   };

var clone1;
var clone2;
var samePerson;

// Create copies
clone1 = Object.assign({}, person);
clone2 = JSON.parse(JSON.stringify(person));
// clone3 = structuredClone(person);
samePerson = person;

// Modify the ORIGINAL object
person.age = 78;
person.country = "FR";

// console.log("clone1: ",clone1);
// console.log("clone2: ",clone2);
// console.log("clone3: ",clone3)
// console.log("samePerson: ",samePerson);
