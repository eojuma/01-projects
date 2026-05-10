const dataSet = `
qqq http:// qqqq q qqqqq https://something.com/hello qqqqqqq 
qhttp://example.com/hello?you=something&something=you
https://page.com/test?x=1&y=2&z=3
https://data.com/path?a=1&b=2&c=3&d=4
`;


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
