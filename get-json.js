function getJSON(path, params) {
  let url = path;
  if (params && Object.keys(params).length > 0) {
    const queryString = new URLSearchParams(params).toString();
    url += (path.includes('?') ? '&' : '?') + queryString;
  }

  return fetch(url)
    .then((response) => {
      if (!response.ok) {
        throw new Error(response.statusText || `HTTP error! status: ${response.status}`);
      }
      return response.json();
    })
    .then((json) => {
      if (json && typeof json === 'object') {
        if ('error' in json) {
          throw new Error(json.error);
        }
        if ('data' in json) {
          return json.data;
        }
      }
      return json;
    });
  }