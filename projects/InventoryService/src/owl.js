const fetch = require("node-fetch");

const brickowl = `https://api.brickowl.com/v1`;
const key = process.env.BRICKOWL_TOKEN;

class Owl {
  constructor(client) {
    this.fetch = client || fetch;
  }

  getBoid(setId) {
    const uri = `${brickowl}/catalog/id_lookup?key=${key}&type=Set&id=${setId}`;

    return fetch(uri)
      .then(res => res.json())
      .then(json => (json.boids.length > 0 ? json.boids[0] : undefined));
  }

  getInventory(boid) {
    const uri = `${brickowl}/catalog/inventory?key=${key}&boid=${boid}`;

    return fetch(uri)
      .then(res => res.json())
      .then(json => {
        if (!json.inventory) {
          return undefined;
        }

        return json.inventory.map(item => {
          return { quantity: Number(item.quantity), boids: [item.boid] };
        });
      });
  }
}

export default Owl;
