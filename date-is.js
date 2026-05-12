function isValid(date) {
    // 1. If it's a Date object, make sure it's not "Invalid Date"
    if (date instanceof Date) {
        return !isNaN(date.getTime());
    }
    
    // 2. If it's a Number, it's valid (Date.now() is a number)
    if (typeof date === 'number') {
        return true;
    }

    // 3. For everything else (strings, booleans, null), return false
    return false;
}

function isAfter(date1, date2) {
    if (!isValid(date1) || !isValid(date2)) return false;
    return new Date(date1) > new Date(date2);
}

function isBefore(date1, date2) {
    if (!isValid(date1) || !isValid(date2)) return false;
    return new Date(date1) < new Date(date2);
}

function isFuture(date) {
    if (!isValid(date)) return false;
    return new Date(date) > new Date();
}

function isPast(date) {
    if (!isValid(date)) return false;
    return new Date(date) < new Date();
}