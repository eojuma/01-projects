function citiesOnly(weatherInfo){
return weatherInfo.map(({city})=> city);
};

function upperCasingStates(states){
    return states.map(state=>state.split(' ').map(names=>names[0].toUpperCase()+names.slice(1).tolowerCase()).join(' '));
}

const fahrenheitToCelsius = (temps) => {
    return temps.map(temp => {
      const f = parseInt(temp);
      const c = Math.floor((f - 32) * 5 / 9);
      return `${c}°C`;
    });
  };

  function trimTemp(temps){
    return temps.map(temp=>({...temp,temperature: temp.temperature.replace(/\s+/g,'')}));
  };

  