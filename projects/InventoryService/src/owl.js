const fetch = require("node-fetch");

const brickowl = `https://api.brickowl.com/v1`;
const key = process.env.BRICKOWL_TOKEN;

const defaultClient = uri => fetch(uri).then(res => res.json());

class Owl {
  constructor(client) {
    this.fetch = client || defaultClient;
  }

  getBoid(setId) {
    const uri = `${brickowl}/catalog/id_lookup?key=${key}&type=Set&id=${setId}`;

    return this.fetch(uri).then(
      json => (json.boids.length > 0 ? json.boids[0] : undefined)
    );
  }

  getInventory(boid) {
    const uri = `${brickowl}/catalog/inventory?key=${key}&boid=${boid}`;

    return this.fetch(uri).then(json => {
      if (!json.inventory) {
        return undefined;
      }

      const grouped = json.inventory.reduce((acc, item) => {
        const id = item.alt_link === 0 ? item.boid : item.alt_link;
        const group = (acc[id] = acc[id] || {
          boids: []
        });

        group.quantity = Number(item.quantity);
        group.boids.push(item.boid);

        return acc;
      }, {});

      return Object.values(grouped).sort((x, y) => x.boids[0] > y.boids[0]);
    });
  }
}

export default Owl;
