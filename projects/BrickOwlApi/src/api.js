import Client from "./client";
import BoidCache from "./boidCache";

const withoutBoids = model => {
  const { boids, ...part } = model;
  return part;
};

const expression = /^(.*)-(.*)$/;
const colorFromBoid = boid => {
  const match = boid.match(expression);

  return match ? Number(match[2]) : undefined;
};
const boidWithoutColor = boid => {
  const match = boid.match(expression);

  return match ? match[1] : boid;
};

const buildPart = (lookup, model) => {
  const boid = model.boids[0];

  return Object.assign(withoutBoids(model), {
    partNumber: lookup[boidWithoutColor(boid)],
    color: colorFromBoid(boid)
  });
};

const lookupInventory = (client, boidCache, setNumber) => {
  return client.getSetBoid(setNumber).then(setBoid => {
    if (!setBoid) {
      return Promise.resolve([]);
    }

    return client
      .getInventory(setBoid)
      .then(inventory =>
        boidCache
          .getMany(inventory.map(part => part.boids[0].replace(/(-.*)/g, "")))
          .then(lookup => inventory.map(model => buildPart(lookup, model)))
      );
  });
};

class BrickOwlApi {
  constructor({ brickOwlToken, storage, client }) {
    const httpClient = client || new Client(brickOwlToken);
    const boidCache = new BoidCache(storage, httpClient);

    this.getInventory = setNumber =>
      lookupInventory(httpClient, boidCache, setNumber);
  }
}

export default BrickOwlApi;
