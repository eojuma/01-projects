// --- Core Utilities (Your existing correct code) ---

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

// --- Corrected Derived Functions ---

// 1. totalCalories: Sums up the 'calories' property from every item in the cart
const totalCalories = (cart) => {
  return reduceEntries(cart, (total, [item, nutrition]) => {
    return total + (nutrition.calories || 0);
  }, 0);
};

const lowCarbs = (cart) => {
  return filterEntries(cart, ([item, nutrition]) => {
    return (nutrition.carbs || 0) < 50;
  });
};


const cartTotal = (cart) => {
  return reduceEntries(cart, (acc, [item, nutrition]) => {
    Object.entries(nutrition).forEach(([macro, value]) => {
      acc[macro] = (acc[macro] || 0) + value;
      
      acc[macro] = Math.round(acc[macro] * 1000) / 1000;
    });
    return acc;
  }, {});
};

export { filterEntries, mapEntries, reduceEntries, totalCalories, lowCarbs, cartTotal };