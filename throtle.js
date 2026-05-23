function throttle(func, delay) {
  let isThrottled = false;

  return function (...args) {
    if (isThrottled) {
      return;
    }
    func.apply(this, args);

    isThrottled = true;
    setTimeout(() => {
      isThrottled = false;
    }, delay);
  };
}