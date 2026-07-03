function pronoun(str) {
  const targetPronouns = new Set(['i', 'you', 'he', 'she', 'it', 'they', 'we']);
  const result = {};

  const words = str.toLowerCase().match(/\b[a-z']+\b/g) || [];

  for (let i = 0; i < words.length; i++) {
    const currentWord = words[i];

    if (targetPronouns.has(currentWord)) {
      if (!result[currentWord]) {
        result[currentWord] = { word: [], count: 0 };
      }

      result[currentWord].count += 1;

      const nextWord = words[i + 1];

      if (nextWord && !targetPronouns.has(nextWord)) {
        result[currentWord].word.push(nextWord);
      }
    }
  }

  return result;
}