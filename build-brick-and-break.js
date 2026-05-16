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

export function repair(...ids){
ids.forEach(id => {
    const brick = document.getElementById(id);
    if (brick){
        if (brick.dataset.foundation === "true"){
            brick.dataset.repaired = "in progress"
        }else{
            brick.dataset.repaired="true"
        }
    }
});
}
export function destroy(){

}