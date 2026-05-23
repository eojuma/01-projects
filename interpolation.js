function interpolation({ step, start, end, duration, callback }) {
  const interval = duration / step;
  const range = end - start;

  const run = async () => {
    for (let i = 0; i < step; i++) {
      const x = i / step;
      
      const y = start + x * range;

      callback([x, y]);

      await delay(interval);
    }
  };

  function delay(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  run();
}