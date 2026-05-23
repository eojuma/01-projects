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
  let lastExecutedTime = 0;

  let leading = true;
  let trailing = false;

  if ('leading' in options || 'trailing' in options) {
    leading = !!options.leading;
    trailing = !!options.trailing;
  }

  const timerExpired = () => {
    const now = Date.now();
    const timeSinceLastExec = now - lastExecutedTime;
    const remaining = wait - timeSinceLastExec;

    if (trailing && savedArgs && remaining <= 0) {
      lastExecutedTime = now;
      const args = savedArgs;
      const thisArg = savedThis;
      savedArgs = null;
      savedThis = null;

      func.apply(thisArg, args);
      timeoutId = setTimeout(timerExpired, wait);
    } else {
      timeoutId = null;
      savedArgs = null;
      savedThis = null;
    }
  };

  return function (...args) {
    const now = Date.now();

    if (!timeoutId) {
      if (leading) {
        lastExecutedTime = now;
        func.apply(this, args);
      } else {
        savedArgs = args;
        savedThis = this;
      }
      timeoutId = setTimeout(timerExpired, wait);
    } else {
      if (trailing) {
        savedArgs = args;
        savedThis = this;
      }
    }
  };
}