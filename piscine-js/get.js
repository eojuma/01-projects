function get(src, path) {
  if (path === "") return undefined;

  let keys = path.split(".");

  for (let i = 0; i < keys.length; i++) {
    if (src === null || src === undefined) {
      return undefined;
    }

    src = src[keys[i]];
  }

  return src;
}
