import { NUTRITION_DB } from './manipulate-entries.data.js';


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

const totalCalories = (cart) => {
  const total = reduceEntries(cart, (sum, [item, weight]) => {
    const caloriesPer100g = NUTRITION_DB[item]?.calories || 0;
    return sum + (caloriesPer100g * (weight / 100));
  }, 0);
  return Math.round(total * 100) / 100; 
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
    const itemMacros = NUTRITION_DB[item] || {};
    
    return [
      item,
      mapEntries(itemMacros, ([macro, valuePer100g]) => {
        const totalValue = valuePer100g * (weight / 100);
        return [macro, Math.round(totalValue * 1000) / 1000]; 
      })
    ];
  });
};

export { filterEntries, mapEntries, reduceEntries, totalCalories, lowCarbs, cartTotal };