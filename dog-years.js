function dogYears(planet,num){
    const secInYear = 31557600;
    const dogY = 7
    var toLowerPlanet = planet.toLowerCase();
       if (planet==="earth"){
       return Number((num/secInYear*dogY).toFixed(2));
    }else if(planet=="mercury"){
        return Number((num/secInYear*dogY/0.2408467).toFixed(2))
    }else if (planet=="venus"){
        return Number((num/secInYear*dogY/0.61519726).toFixed(2))

    }else if (planet=="mars"){
                return Number((num/secInYear*dogY/1.8808158).toFixed(2))

    }else if (planet=="jupiter"){
                return Number((num/secInYear*dogY/11.862615).toFixed(2))

    }else if (planet=="saturn"){
                return Number((num/secInYear*dogY/29.447498).toFixed(2))

    }else if (planet=="uranus"){
                return Number((num/secInYear*dogY/84.016846).toFixed(2))

    }else if (planet=="neptune"){
                return Number((num/secInYear*dogY/164.79132).toFixed(2))

    }

}

// console.log((dogYears("earth", 1000000000))); // Should output 221.82
// console.log(dogYears("mercury", 1000000000)); // Should output 921.3
// console.log(dogYears("venus", 1000000000)); // Should output 360.59