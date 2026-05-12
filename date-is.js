function isValid(date) {
    // 1. Convert the input to a Date object first
    const d = new Date(date);
    
    // 2. Use the new variable 'd' to check validity
    // We check if d is a Date and its value isn't NaN
    // Important: we check 'typeof date !== "boolean"' because 
    // new Date(true) is technically valid in JS but usually shouldn't be.
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