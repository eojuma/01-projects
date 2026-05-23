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

  const leading = options.leading !== false;
  const trailing = options.trailing !== false;

  const invokeFunc = (time) => {
    lastExecutedTime = time;
    const args = savedArgs;
    const thisArg = savedThis;
    savedArgs = null;
    savedThis = null;
    func.apply(thisArg, args);
  };

  const startTimer = (pendingWait) => {
    timeoutId = setTimeout(timerExpired, pendingWait);
  };

  const timerExpired = () => {
    timeoutId = null;
    const now = Date.now();
    const timeSinceLastExec = now - lastExecutedTime;
    const remaining = wait - timeSinceLastExec;

    if (trailing && savedArgs) {
      if (remaining <= 0) {
        invokeFunc(now);
        startTimer(wait);
      } else {
        startTimer(remaining);
      }
    }
  };

  return function (...args) {
    const now = Date.now();
    const isFirstCall = !timeoutId;

    if (isFirstCall) {
      lastExecutedTime = now;
      if (leading) {
        func.apply(this, args);
      } else {
        savedArgs = args;
        savedThis = this;
      }
      startTimer(wait);
    } else {
      savedArgs = args;
      savedThis = this;
    }
  };
}