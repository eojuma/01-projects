function queryServers(serverName, q) {
  const urlPrimary = `/${serverName}?q=${q}`;
  const urlBackup = `/${serverName}_backup?q=${q}`;

  return Promise.race([
    getJSON(urlPrimary),
    getJSON(urlBackup)
  ]);
}

function gougleSearch(q) {
  const webQuery = queryServers('web', q);
  const imageQuery = queryServers('image', q);
  const videoQuery = queryServers('video', q);

  const allQueriesPromise = Promise.all([webQuery, imageQuery, videoQuery])
    .then(([webResult, imageResult, videoResult]) => {
      return {
        web: webResult,
        image: imageResult,
        video: videoResult
      };
    });

  const timeoutPromise = new Promise((_, reject) => {
    setTimeout(() => {
      reject(new Error('timeout'));
    }, 80);
  });

  return Promise.race([allQueriesPromise, timeoutPromise]);
}