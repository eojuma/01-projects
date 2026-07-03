//when sunday is set to be the first day of the week(and 0 based indexing)
//the friday is in the sixth position but index five

function isFriday(date) {
  if (!(date instanceof Date) || isNaN(date)) return false;
  return date.getDay() === 5;
}

function isWeekend(date) {
  if (!(date instanceof Date) || isNaN(date)) return false;
  const day = date.getDay();
  return day === 0 || day === 6;
  // in this case 0 is on Sunday and 0 is on Saturday
}

function isLeapYear(date) {
  if (!(date instanceof Date) || isNaN(date)) return false;
  const year = date.getFullYear();

  return (year % 4 === 0 && year % 100 !== 0) || year % 400 === 0;
}

function isLastDayOfMonth(date) {
  if (!(date instanceof Date) || isNaN(date)) return false;

  const year = date.getFullYear();
  const month = date.getMonth();
  const day = date.getDate();

  const lastDay = new Date(year, month + 1, 0).getDate();

  return day === lastDay;
}

console.log(isFriday(new Date("2024-03-01"))); // true
console.log(isWeekend(new Date("2024-03-03"))); // true
console.log(isLeapYear(new Date("2024-01-01"))); // true
console.log(isLastDayOfMonth(new Date("2024-02-29"))); // true
console.log(isLastDayOfMonth(new Date("2024-02-28"))); // false