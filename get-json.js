async function getJSON(path, params) {
  let url = path;
  if (params && Object.keys(params).length > 0) {
    let queryString = new URLSearchParams(params).toString();
    
    queryString = queryString.replace(/\+/g, '%20');
    
    url += (path.includes('?') ? '&' : '?') + queryString;
  }

  const response = await fetch(url);

  if (response && typeof response === 'object' && 'ok' in response && response.ok === false) {
    throw new Error(response.statusText || 'HTTP error');
  }

  const json = response && typeof response.json === 'function' 
    ? await response.json() 
    : response;

  if (json && typeof json === 'object') {
    if ('error' in json) {
      throw new Error(json.error);
    }
    if ('data' in json) {
      return json.data;
    }
  }

  return json;
}