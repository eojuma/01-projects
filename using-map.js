function citiesOnly(weatherInfo){
return weatherInfo.map(({city})=> city);
};

function upperCasingStates(states){
    return states.map(state=>state.split(' ').map(names=>names[0].toUpperCase()+names.slice(1).toLowerCase()).join(' '));
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

  function tempForecasts(weatherInfo){
    return weatherInfo.map(item=>{
        const fahrenheitValue = parseInt(item.temperature.replace(/\s+/g, ''));
        const celciusValue = Math.floor((fahrenheitValue-32)*5/9);

        const formattedState = item.state
        .split(' ')
        .map(word => word[0].toUpperCase() + word.slice(1).toLowerCase())
        .join(' ');
        return `${cValue}°Celsius in ${item.city}, ${formattedState}`;    });
  }