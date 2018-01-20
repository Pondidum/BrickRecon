import Client from "./client";
import BoidCache from "./boidCache";

const withoutBoids = model => {
  const { boids, ...part } = model;
  return part;
};

const colorFromBoid = boid => Number(boid.replace(/.*-(.*)$/, "$1"));

const buildPart = (lookup, model) => {
  const boid = model.boids[0];

  return Object.assign(withoutBoids(model), {
    partNumber: lookup[boid],
    color: colorFromBoid(boid)
  });
};

const lookupInventory = (client, boidCache, setNumber) => {
  return client
    .boidFromSetNumber(setNumber)
    .then(setBoid => client.getInventory(setBoid))
    .then(inventory =>
      boidCache
        .getMany(inventory.map(part => part.boids[0]))
        .then(lookup => inventory.map(model => buildPart(lookup, model)))
    );
};

class BrickOwlApi {
  constructor({ brickOwlToken, storage, client }) {
    const httpClient = client || new Client(brickOwlToken);
    const boidCache = new BoidCache(storage, client);

    this.getInventory = setNumber =>
      lookupInventory(httpClient, boidCache, setNumber);
  }
}

export default BrickOwlApi;
