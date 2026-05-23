function interpolation({ step, start, end, duration, callback }) {
  const intervalTime = duration / step;

  for (let i = 0; i < step; i++) {
    const x = i / step;
    const y = (i + 1) * intervalTime;
    setTimeout(() => {
      callback([x, y]);
    },y);
  }
}