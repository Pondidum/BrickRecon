import { mapFrom } from "./util";

const get = (store, boids) => {
  const map = mapFrom(
    boids.filter(boid => store[boid]),
    boid => boid,
    boid => store[boid]
  );

  return Promise.resolve(map);
};

const write = (store, boids) => {
  Object.assign(store, boids);
  return Promise.resolve();
};

class MemoryStorage {
  constructor() {
    const store = {};
    this.getMany = boids => get(store, boids);
    this.writeMany = boids => write(store, boids);
  }
}

export default MemoryStorage;
