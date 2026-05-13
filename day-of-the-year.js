function dayOfTheYear(date) {
  const year = date.getUTCFullYear();

  const startOfYear = Date.UTC(year, 0, 1);

  const currentTimestamp = date.getTime();

  const diffInMilliseconds = currentTimestamp - startOfYear;

  const msPerDay = 1000 * 60 * 60 * 24;
  return Math.floor(diffInMilliseconds / msPerDay) + 1;
}
