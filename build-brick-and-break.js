export function build(totalbrics){

let currentBrick  = 1;
const intervals = setInterval(()=>{
    if (currentBrick>totalbrics){
        clearInterval(intervals);
        return;
    }
    const brick = document.createElement('div');
    brick.id = `brick-${currentBrick}`;
    if (currentBrick%3==2){
brick.dataset.foundation ="true";
    }
    document.body.appendChild(brick);
    currentBrick++;
},100);

}
