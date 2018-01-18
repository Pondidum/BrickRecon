export const chunk = (arr, size) => {
  var results = [];

  while (arr.length) {
    results.push(arr.splice(0, size));
  }

  return results;
};

export const mapFrom = (arr, keyFunc, valueFunc = value => value) =>
  arr.reduce((all, current) => {
    const key = keyFunc(current);
    const value = valueFunc(current);

    all[key] = value;
    return all;
  }, {});
