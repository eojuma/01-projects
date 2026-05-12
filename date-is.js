function isValid(date) {
    // 1. Convert the input to a Date object first
    const d = new Date(date);
    
    return d instanceof Date && !isNaN(d) && typeof date !== 'boolean';
}

function isAfter(date1, date2) {
    // Both must be valid dates first
    if (!isValid(date1) || !isValid(date2)) return false;
    return new Date(date1) > new Date(date2);
}

function isBefore(date1, date2) {
    if (!isValid(date1) || !isValid(date2)) return false;
    return new Date(date1) < new Date(date2);
}

function isFuture(date) {
    const d = new Date(date);
    // Use the converted date 'd' for all checks
    return isValid(date) && d > new Date();
}

function isPast(date) {
    const d = new Date(date);
    // Use the converted date 'd' for all checks
    return isValid(date) && d < new Date();
}