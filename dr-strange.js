const DayNames = [
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
    "Sunday",
    "SecondMonday",
    "SecondTuesday",
    "SecondWednesday",
    "SecondThursday",
    "SecondFriday",
    "SecondSaturday",
    "SecodSunday",

]


function addWeek(date){

    const epoch = new Date(1, 0, 1); 
    epoch.setFullYear(1); 

    const diffInMs = date.getTime() - epoch.getTime();
    const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

  
    const cycleIndex = ((diffInDays % 14) + 14) % 14;

    return DayNames[cycleIndex];
}


function timeTravel({ date, hour, minute, second }) {
    const newDate = new Date(date.getTime());
    newDate.setHours(hour, minute, second, 0);
    return newDate;
}