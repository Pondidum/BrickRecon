const chunk = (arr, size) => {
  var results = [];

  while (arr.length) {
    results.push(arr.splice(0, size));
  }

  return results;
};

const client = {
  boidFromSetNumber: setNumber => Promise.resolve(123123),
  getInventory: boid => Promise.resolve([]),
  getPartNumbers: boids => Promise.resolve({})
};

const storage = {
  getMany: boids => Promise.resolve({}),
  writeMany: boids => Promise.resolve()
};

const boidCache = {
  getMany: boids =>
    storage
      .getMany(boids) //make this handle 100 items max internally
      .then(map =>
        client
          .getPartNumbers(Object.keys(map).filter(boid => !map[boid])) //make this handle 100 items max internally
          .then(partNumbers =>
            storage
              .writeMany(partNumbers)
              .then(() => Object.assign({}, map, partNumbers))
          )
      )
};

const lookupInventory = setNumber => {
  const boid = client.boidFromSetNumber(setNumber);
  const inventory = client.getInventory(boid);

  return boidCache.getMany(inventory.map(part => part.boid)).then(partNumbers =>
    inventory.map(model => {
      const { boid, ...part } = model;
      return Object.assign({}, part, { partNumber: partNumbers[part.boid] });
    })
  );
};
