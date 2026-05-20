
const filterEntries = (obj, callback) => {
  return Object.fromEntries(
    Object.entries(obj).filter(([key, value]) => callback([key, value]))
  );
};

const mapEntries = (obj, callback) => {
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => callback([key, value]))
  );
};

const reduceEntries = (obj, callback, initialValue) => {
  return Object.entries(obj).reduce(
    (acc, [key, value]) => callback(acc, [key, value]),
    initialValue
  );
};


const round3 = (num) => Math.round(num * 1000) / 1000;

const totalCalories = (cart) => {
  return reduceEntries(cart, (total, [item, weight]) => {
    const caloriesPer100g = NUTRITION_DB[item]?.calories || 0;
    return total + (caloriesPer100g * (weight / 100));
  }, 0);
};

const lowCarbs = (cart) => {
  return filterEntries(cart, ([item, weight]) => {
    const carbsPer100g = NUTRITION_DB[item]?.carbs || 0;
    const totalCarbsForItem = carbsPer100g * (weight / 100);
    return totalCarbsForItem < 50;
  });
};


const cartTotal = (cart) => {
  return mapEntries(cart, ([item, weight]) => {
    const macrosPer100g = NUTRITION_DB[item] || {};
    
    const calculatedNutrition = mapEntries(macrosPer100g, ([macro, valuePer100g]) => {
      return [macro, round3(valuePer100g * (weight / 100))];
    });

    return [item, calculatedNutrition];
  });
};