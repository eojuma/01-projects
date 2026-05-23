function interpolation({ step, start, end, duration, callback }) {
  const intervalTime = duration / step;
  const range = end - start;

  for (let i = 0; i < step; i++) {
    const x = i / step;
    const y = start + x * range;

    const executionDelay = (i + 1) * intervalTime;

    setTimeout(() => {
      callback([x, y]);
    }, executionDelay);
  }
}