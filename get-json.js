async function getJSON(path, params) {
  let url = path;
  if (params && Object.keys(params).length > 0) {
    const queryString = new URLSearchParams(params).toString();
    url += (path.includes('?') ? '&' : '?') + queryString;
  }

  const response = await fetch(url);

  if (response && response.ok === false) {
    throw new Error(response.statusText || 'HTTP error');
  }

  let json;
  try {
    json = response && typeof response.json === 'function' 
      ? await response.json() 
      : response;
  } catch (e) {
    json = response;
  }

  if (json && typeof json === 'object') {
    if ('error' in json && typeof json.error === 'string') {
      throw new Error(json.error);
    }
    if ('data' in json) {
      return json.data;
    }
  }

  return json;
}