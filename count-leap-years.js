
function countLeapYears(date) {
    const year = date.getFullYear() - 1;

    if (year < 4) return 0; 

    const count4 = Math.floor(year / 4);
    const count100 = Math.floor(year / 100);
    const count400 = Math.floor(year / 400);

    return count4 - count100 + count400;
}