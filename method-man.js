function words(s) {
  return s.split(" ");
}

function sentence(s) {
  const together = s.join(" ");
  return together;
}

function yell(s) {
  const upper = s.toUpperCase();
  return upper;
}

function whisper(s) {
 const lower = s.toLowerCase();
  return `*${lower}*`;
}

function capitalize(s) {
  const transformed = s.charAt(0).toUpperCase() + s.slice(1).toLowerCase();
  return transformed;
}
