function interpolation({ step, start, end, duration, callback }) {
  const interval = duration / step;
  const range = end - start;

  function runStep(i) {
    if (i >= step) return;

    const x = i / step;
    const y = start + x * range;

    callback([x, y]);

    setTimeout(() => {
      runStep(i + 1);
    }, interval);
  }

  runStep(0);
}