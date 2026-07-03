function adder(numbers,start=0) {
    return numbers.reduce((accumulator, current) => accumulator + current, start);
  }

  function sumOrMul(numbers, initialValue = 0) {
    return numbers.reduce((accumulator, current) => {
      if (current % 2 !== 0) {
        return accumulator + current;
      } else {
        return accumulator * current;
      }
    }, initialValue);
  }

  function funcExec(functions, initialValue) {
    return functions.reduce((accumulator, currentFunc) => {
      return currentFunc(accumulator);
    }, initialValue);
  }