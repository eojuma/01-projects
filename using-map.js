function citiesOnly(weatherInfo){
return weatherInfo.map(({city})=> city);
};

function upperCasingStates(states){
    return states.map(state=>state.split(' ').map(names=>names[0].toUpperCase()+names.slice(1).tolowerCase()).join(' '));}