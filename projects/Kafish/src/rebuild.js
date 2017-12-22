const convert = (key, value) => {
  if (key === "N") {
    return Number.parseFloat(value);
  }

  return value;
};

export default record =>
  Object.keys(record).reduce((out, prop) => {
    const element = record[prop];
    const key = Object.keys(element)[0];

    out[prop] = convert(key, element[key]);
    return out;
  }, {});
