function all(obj) {
  return new Promise((resolve, reject) => {
    if (!obj || typeof obj !== 'object') {
      return resolve({});
    }

    const entries = Object.entries(obj);
    const totalCount = entries.length;

    if (totalCount === 0) {
      return resolve({});
    }

    const resultObject = {};
    let resolvedCount = 0;
    let hasFailed = false;

    entries.forEach(([key, value]) => {
      Promise.resolve(value)
        .then((resolvedValue) => {
          if (hasFailed) return;

          resultObject[key] = resolvedValue;
          resolvedCount++;

          if (resolvedCount === totalCount) {
            resolve(resultObject);
          }
        })
        .catch((error) => {
          if (!hasFailed) {
            hasFailed = true;
            reject(error);
          }
        });
    });
  });
}