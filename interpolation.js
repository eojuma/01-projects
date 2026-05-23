function interpolation({ step, start, end, duration, callback }) {
  const intervalTime = duration / step;
  const range = end - start;
  let i = 0;

  const intervalId = setInterval(() => {
    if (i >= step) {
      clearInterval(intervalId);
      return;
    }

    const x = i / step;
    const y = start + x * range;

    callback([x, y]);
    i++;
  }, intervalTime);

  const x = 0;
  const y = start;
  callback([x, y]);
  i++;
}