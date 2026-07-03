function debounce(func, wait) {
  let timeoutId = null;

  return function (...args) {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    timeoutId = setTimeout(() => {
      func.apply(this, args);
    }, wait);
  };
}


function opDebounce(func, wait, options = {}) {
  let timeoutId = null;
  const leading = !!options.leading;

  return function (...args) {
    const shouldCallNow = leading && !timeoutId;

    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    timeoutId = setTimeout(() => {
      timeoutId = null;

      if (!leading) {
        func.apply(this, args);
      }
    }, wait);

    if (shouldCallNow) {
      func.apply(this, args);
    }
  };
}