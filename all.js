function all(obj) {
  if (!obj || typeof obj !== 'object') {
    return Promise.resolve({});
  }

  const entries = Object.entries(obj);

  const promises = entries.map(([key, value]) => {
    return Promise.resolve(value).then((resolvedValue) => {
      return [key, resolvedValue];
    });
  });

  return Promise.all(promises).then((resolvedEntries) => {
    return Object.fromEntries(resolvedEntries);
  });
}