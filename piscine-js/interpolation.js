function interpolation({ step, start, end, duration, callback }) {
  const intervalTime = duration / step;
  const range = end - start;

  for (let i = 0; i < step; i++) {
    const progress = i / step;

    let val = start + progress * range;

    let time = (i + 1) * intervalTime;

    const round = (num) => {
      return Number.isInteger(num) ? num : Math.round(num * 100) / 100;
    };

    const finalVal = round(val);
    const finalTime = round(time);

    setTimeout(() => {
      callback([finalVal, finalTime]);
    }, time);
  }
}