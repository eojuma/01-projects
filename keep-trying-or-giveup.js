function retry(count, callback) {
  return async function (...args) {
    const maxAttempts = 1 + count;

    for (let attempt = 1; attempt <= maxAttempts; attempt++) {
      try {
        return await callback.apply(this, args);
      } catch (error) {
        if (attempt === maxAttempts) {
          throw new Error("All retries failed");
        }
      }
    }
  };
}

function timeout(delay, callback) {
  return function (...args) {
    const actionPromise = callback.apply(this, args);

    const timeoutPromise = new Promise((_, reject) => {
      setTimeout(() => {
        reject(new Error("timeout"));
      }, delay);
    });

    return Promise.race([actionPromise, timeoutPromise]);
  };
}