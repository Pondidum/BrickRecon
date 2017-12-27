class OwlDecorator {
  constructor(cache, owl) {
    this.cache = cache;
    this.owl = owl;
  }

  getInventory(boid) {
    const partsFromBoids = item =>
      Promise.all(item.boids.map(boid => this.cache.get(boid)));

    const buildPart = (item, parts) => ({
      partNumber: parts[0].partNumber,
      color: parts[0].color,
      quantity: item.quantity,
      otherPartNumbers: parts.slice(1).map(part => part.partNumber)
    });

    const replaceBoids = item =>
      partsFromBoids(item).then(parts => buildPart(item, parts));

    return this.owl
      .getInventory(boid)
      .then(inv => inv.map(item => replaceBoids(item)))
      .then(promises => Promise.all(promises));
  }
}

export default OwlDecorator;
