function isFriday(date) {
  if (!(date instanceof Date) || isNaN(date)) return false;
  return date.getDay() === 5;
}

function isWeekend(date) {
  if (!(date instanceof Date) || isNaN(date)) return false;
  const day = date.getDay();
  return day === 0 || day === 6;
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
