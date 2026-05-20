const nutritionDB = {
  tomato: {
    calories: 18,
    protein: 0.9,
    carbs: 3.9,
    sugar: 2.6,
    fiber: 1.2,
    fat: 0.2,
  },
  vinegar: {
    calories: 20,
    protein: 0.04,
    carbs: 0.6,
    sugar: 0.4,
    fiber: 0,
    fat: 0,
  },
  oil: { calories: 48, protein: 0, carbs: 0, sugar: 123, fiber: 0, fat: 151 },
  onion: { calories: 0, protein: 1, carbs: 9, sugar: 0, fiber: 0, fat: 0 },
  garlic: {
    calories: 149,
    protein: 6.4,
    carbs: 33,
    sugar: 1,
    fiber: 2.1,
    fat: 0.5,
  },
  paprika: {
    calories: 282,
    protein: 14.14,
    carbs: 53.99,
    sugar: 1,
    fiber: 0,
    fat: 12.89,
  },
  sugar: {
    calories: 387,
    protein: 0,
    carbs: 100,
    sugar: 100,
    fiber: 0,
    fat: 0,
  },
  orange: {
    calories: 49,
    protein: 0.9,
    carbs: 13,
    sugar: 9,
    fiber: 0.2,
    fat: 0.1,
  },
};

function filterEntries(obj, predicate) {
  return Object.fromEntries(Object.entries(obj).filter(predicate));
}

function mapEntries(obj, transform) {
  return Object.fromEntries(Object.entries(obj).map(transform));
}

function reduceEntries(obj, reducer, initialValue) {
  return Object.entries(obj).reduce(reducer, initialValue);
}
function totalCalories(groceriesCart) {
  return Number(
    reduceEntries(
      groceriesCart,
      (total, [item, grams]) => {
        return total + (nutritionDB[item].calories * grams) / 100;
      },
      0,
    ).toFixed(1),
  );
}

function lowCarbs(groceriesCart) {
  return filterEntries(groceriesCart, ([item, grams]) => {
    return (nutritionDB[item].carbs * grams) / 100 < 50;
  });
}

function cartTotal(groceriesCart) {
  return mapEntries(groceriesCart, ([item, grams]) => {
    const per100g = nutritionDB[item];
    const nutrition = {};
    for (let nutrient in per100g) {
      let value = (per100g[nutrient] * grams) / 100;
      nutrition[nutrient] = parseFloat(value.toFixed(3));
    }
    return [item, nutrition];
  });
}
