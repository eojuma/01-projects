
function getURL(dataSet){
    const regex = /https?:\/\/[^\s]+/g;
    return dataSet.match(regex) || [];
}


function greedyQuery(dataSet){
    const allurls=dataSet.match(/https?:\/\/[^\s]+/g) || [];
    const triple=/\?.*&.*&.*/;
    return allurls.filter(url=>triple.test(url));
}


function notSoGreedy(text) {
  const allUrls = text.match(/https?:\/\/[^\s]+/g) || [];

  return allUrls.filter(url => {
    if (!url.includes('?')) return false;

    const ampersandCount = (url.split('&').length - 1);

    return ampersandCount >= 1 && ampersandCount <= 2;
  });
}
