async function series(tasks) {
  const container = [];
  for (const task of tasks) {
    container.push(await task());
  }
  return container;
}