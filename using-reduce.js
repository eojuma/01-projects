function adder(numbers) {
    return numbers.reduce((accumulator, current) => accumulator + current, 0);
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