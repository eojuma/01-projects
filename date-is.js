function isValid(date) {
    // 1. Strict Type Check: If it's not an instance of Date, return false immediately
    if (!(date instanceof Date)) {
        return false;
    }
    // 2. Value Check: If it is a Date object, make sure it isn't "Invalid Date" (NaN)
    return !isNaN(date.getTime());
}

function isAfter(date1, date2) {
    // Both must be valid Date objects according to our strict isValid
    if (!isValid(date1) || !isValid(date2)) return false;
    return date1 > date2;
}

function isBefore(date1, date2) {
    if (!isValid(date1) || !isValid(date2)) return false;
    return date1 < date2;
}

function isFuture(date) {
    // Must be a valid Date object AND in the future
    return isValid(date) && date > new Date();
}

function isPast(date) {
    // Must be a valid Date object AND in the past
    return isValid(date) && date < new Date();
}