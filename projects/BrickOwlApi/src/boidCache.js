const getMany = (storage, client, boids) =>
  storage
    .getMany(boids) //make this handle 100 items max internally
    .then(map => {
      const missing = boids.filter(boid => !map[boid]);

      if (missing.length === 0) {
        return map;
      }

      return client
        .getPartNumbers(missing) //make this handle 100 items max internally
        .then(partNumbers =>
          storage
            .writeMany(partNumbers)
            .then(() => Object.assign({}, map, partNumbers))
        );
    });

class BoidCache {
  constructor(storage, client) {
    this.getMany = boids => getMany(storage, client, boids);
  }
}

export default BoidCache;
