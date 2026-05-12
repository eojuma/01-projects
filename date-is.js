function isValid(date){
    return new Date(date) instanceof Date && !isNaN(date.getTime());
}


function isAfter(date1,date2){
    return new Date(date1)>new Date(date2)
}

function isBefore(date1,date2){
    return new Date(date1)<new Date(date2)
}


function isFuture(date){
        return new Date(date) instanceof Date && !isNaN(date.getTime()) && new Date(date) > new Date();

}

function isPast(date){
            return new Date(date) instanceof Date && !isNaN(date.getTime()) && new Date(date) < new Date();

}
