import Client from "./client";
import BoidCache from "./boidCache";

const storage = {
  getMany: boids => Promise.resolve({}),
  writeMany: boids => Promise.resolve()
};

const client = new Client(process.env.BRICKOWL_TOKEN);
const boidCache = new BoidCache(storage, client);

const chunk = (arr, size) => {
  var results = [];

  while (arr.length) {
    results.push(arr.splice(0, size));
  }

  return results;
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
