function dayOfTheYear(date) {
  const year = date.getFullYear();

  const startOfYear = new Date(year, 0, 1);

  const diffInMilliseconds = date - startOfYear;

  const msPerDay = 1000 * 60 * 60 * 24;

  const dayNumber = Math.floor(diffInMilliseconds / msPerDay) + 1;

  return dayNumber;
}
