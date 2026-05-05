const person = {
    name: "Rick",
    age: 77,
    country: "US",
  };
  
  var clone1;
  var clone2;
  var samePerson;
  
  // Create copies
  clone1 = Object.assign({}, person);
  clone2 = JSON.parse(JSON.stringify(person));
  samePerson = person;
  
  // Modify the ORIGINAL object
  person.age = 78;
  person.country = "FR";