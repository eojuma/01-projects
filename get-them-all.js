export function getArchitects(){

    const allElements = Array.from(document.querySelectorAll('#architect-container'));

    const architects = [];
    const nonArchitects = [];

    allElements.forEach(element =>{
       if (element.tagName === 'A'){

           architects.push(element);
       }else{
        nonArchitects.push(element);
       }
    })
    return [architects,nonArchitects];
}

export function getClassical(){
    const architects = Array.from(document.querySelectorAll('a'));
    const classicalArchitects = [];
    const nonClassicalArchitects = [];

    architects.forEach(architect=>{
        if (architect.classList.contains('classical')){
            classicalArchitects.push(architect);
        }else{
            nonClassicalArchitects.push(architect);
        }
    });
    return [classicalArchitects,nonClassicalArchitects];
}

export function getActive(){
    const classicalArchitects = Array.from(document.querySelectorAll('a.classical'));
    const activeClassical = [];
    const inactiveClassical = [];

    classicalArchitects.forEach(architech=>{
        if (architech.classList.contains('active')){
            activeClassical.push(arhitect);
        }else{
            inactiveClassical.push(architect);
        }
    });
    return [activeClassical,inactiveClassical];
}

export function getBonannoPisano(){
    const bonannoPisano = document.getElementById('BonannoPisano');

    const allActiveClassical = Array.from(querySelectorAll('a.classical.active'));
    const otherActiveClassical=allActiveClassical.filter(architech=>architech !==bonannoPisano);
    return [bonannoPisano,otherActiveClassical];
}
