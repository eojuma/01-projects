function hasCity(country, cities) {
    return function(city) {
      // We check if the city exists in the array
      const exists = cities.includes(city);
      
      if (exists) {
        return `${city} is a city from ${country}`;
      } else {
        return `${city} is not a city from ${country}`;
      }
    };
  }