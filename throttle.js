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



function opThrottle(func, wait, options = {}) {
  let timeoutId = null;
  let savedArgs = null;
  let savedThis = null;

  const leading = options.leading !== false;
  const trailing = options.trailing !== false;

  const resetCooldown = () => {
    if (trailing && savedArgs) {
      func.apply(savedThis, savedArgs);
      
      savedArgs = null;
      savedThis = null;
      
      timeoutId = setTimeout(resetCooldown, wait);
    } else {
      timeoutId = null;
    }
  };

  return function (...args) {
    if (!timeoutId) {
      if (leading) {
        func.apply(this, args);
      } else if (trailing) {
        savedArgs = args;
        savedThis = this;
      }

      timeoutId = setTimeout(resetCooldown, wait);
    } else {
      savedArgs = args;
      savedThis = this;
    }
  };
}