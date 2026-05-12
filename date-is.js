function isValid(date) {
    // 1. If it's a Date object, check if it's "Invalid Date"
    if (date instanceof Date) {
        return !isNaN(date.getTime());
    }
    
    // 2. If it's a Number, it MUST NOT be NaN
    // Number.isNaN is the most reliable way to check this
    if (typeof date === 'number' && !Number.isNaN(date)) {
        return true;
    }

    // 3. Everything else (strings, booleans, null, undefined) fails
    return false;
}

// Keep your other functions as they were, ensuring they use the new isValid:
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