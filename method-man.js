function words(s) {
  return s.split(" ");
}

function sentence(s) {
  together = s.join(" ");
  return together;
}

function yell(s) {
  upper = s.toUpperCase();
  return upper;
}

function whisper(s) {
  lower = s.toLowerCase();
  return `*${lower}*`;
}

function capitalize(s) {
  transformed = s.charAt(0).toUpperCase() + s.slice(1).toLowerCase();
  return transformed;
}
