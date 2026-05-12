function firstDayWeek(week, year) {
  let jan1 = new Date(year, 0, 1);

  jan1.setFullYear(parseInt(year));

  let dayOfWeek = jan1.getDay();

  let daysToFirstMonday =
    dayOfWeek === 1 ? 0 : dayOfWeek === 0 ? 1 : 8 - dayOfWeek;

  let targetMonday = new Date(jan1);
  targetMonday.setDate(jan1.getDate() + daysToFirstMonday + (week - 1) * 7);

  if (targetMonday.getFullYear() < parseInt(year)) {
    targetMonday = new Date(jan1);
  }

  let dd = String(targetMonday.getDate()).padStart(2, "0");
  let mm = String(targetMonday.getMonth() + 1).padStart(2, "0");
  let yyyy = String(targetMonday.getFullYear()).padStart(4, "0");

  return `${dd}-${mm}-${yyyy}`;
}
