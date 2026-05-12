function isValid(date){
    return date instanceof Date && !isNaN(date.getTime());
}


function isAfter(date1,date2){
    return date1>date2
}

function isBefore(date1,date2){
    return date1<date2
}


function isFuture(date){
        return date instanceof Date && !isNaN(date.getTime()) && date > new Date();

}

function isPast(date){
            return date instanceof Date && !isNaN(date.getTime()) && date < new Date();

}
