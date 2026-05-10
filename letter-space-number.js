function letterSpaceNumber(s) {
  const pattern = /\w\s\d\b/g;

  const matches = s.match(pattern);

  return matches || [];
}
