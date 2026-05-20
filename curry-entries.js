

const defaultCurry = (obj1) => (obj2) => ({
  ...obj1,
  ...obj2
});

const mapCurry = (fn) => (obj) => {
  const mappedEntries = Object.entries(obj).map(fn);
  return Object.fromEntries(mappedEntries);
};

const reduceCurry = (fn) => (obj, initialValue) => {
  return Object.entries(obj).reduce(fn, initialValue);
};

const filterCurry = (fn) => (obj) => {
  const filteredEntries = Object.entries(obj).filter(fn);
  return Object.fromEntries(filteredEntries);
};



const reduceScore = (obj) => {
  return reduceCurry((acc, [_, person]) => {
    if (person.isForceUser) {
      return acc + person.pilotingScore + person.shootingScore;
    }
    return acc;
  })(obj, 0);
};


const filterForce = (obj) => {
  return filterCurry(([_, person]) => {
    return person.isForceUser && person.shootingScore >= 80;
  })(obj);
};


const mapAverage = (obj) => {
  return mapCurry(([name, person]) => {
    const averageScore = (person.pilotingScore + person.shootingScore) / 2;
    return [
      name,
      { ...person, averageScore }
    ];
  })(obj);
};