function dayOfTheYear(date) {
  const year = date.getUTCFullYear();

  const startOfYear = Date.UTC(year, 0, 1);

  const currentTimestamp = date.getTime();

  const diffInMilliseconds = currentTimestamp - startOfYear;

  const msPerDay = 1000 * 60 * 60 * 24;
  return Math.floor(diffInMilliseconds / msPerDay) + 1;
}
function dayOfTheYear(date) {
    const targetDate = new Date(date.getTime());
    
    const startOfYear = new Date(date.getTime());
    startOfYear.setFullYear(targetDate.getFullYear());
    startOfYear.setMonth(0); 
    startOfYear.setDate(1); 
    
    const diff = targetDate - startOfYear;
    
    const msPerDay = 1000 * 60 * 60 * 24;
    return Math.round(diff / msPerDay) + 1;
}