function sameAmount(s,reg1,reg2){
const match1 = s.match(new RegExp(reg1, 'g')) || [];

const match2 = s.match(new RegExp(reg2, 'g')) || [];

return match1.length === match2.length;


}