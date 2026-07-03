function matchCron(cronString, date) {
    const parts = cronString.split(' ');
    
    const minute = date.getMinutes();
    const hour = date.getHours();
    const dayOfMonth = date.getDate();
    const month = date.getMonth() + 1; 
    
    let dayOfWeek = date.getDay();
    if (dayOfWeek === 0) dayOfWeek = 7;

    const dateValues = [minute, hour, dayOfMonth, month, dayOfWeek];

    for (let i = 0; i < 5; i++) {
        const cronPart = parts[i];
        const dateValue = dateValues[i];

        if (cronPart === '*') {
            continue;
        }

        if (parseInt(cronPart) !== dateValue) {
            return false; 
    }
    }
    return true; 
}