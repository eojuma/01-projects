export function getArchitects(){

    const architects = Array.from(document.querySelectorAll('a'));
const nonArchitects = Array.from(document.querySelectorAll('span'));
    return [architects,nonArchitects];
}

export function getClassical(){
    const classicalArchitects = Array.from(document.querySelectorAll('a.classical'));
   
  const nonClassicalArchitects = Array.from(document.querySelectorAll('a:not(.classical)'));
    return [classicalArchitects,nonClassicalArchitects];
}

export function getActive(){
    const activeClassical = Array.from(document.querySelectorAll('a.classical.active'));
const inactiveClassical = Array.from(document.querySelectorAll('a.classical:not(.active)'));    
    return [activeClassical,inactiveClassical];
}

export function getBonannoPisano(){
    const bonannoPisano = document.getElementById('BonannoPisano');

    const allActiveClassical = Array.from(document.querySelectorAll('a.classical.active'));
    const otherActiveClassical=allActiveClassical.filter(architech=>architech !==bonannoPisano);
    return [bonannoPisano,otherActiveClassical];
}
