import fetch from "node-fetch";
import queryString from "query-string";

const defaultClient = uri => fetch(uri).then(res => res.json());

class Owl {
  constructor(token, client) {
    this.token = token;
    this.fetch = client || defaultClient;
  }

  buildQuery(path, query) {
    const map = Object.assign({ key: this.token }, query);
    const qs = queryString.stringify(map);

    return `https://api.brickowl.com/v1/catalog/${path}?${qs}`;
  }

  getSetBoid(setId) {
    const uri = this.buildQuery("id_lookup", { type: "Set", id: setId });
    const validBoids = js => js && !js.error && js.boids && js.boids.length > 0;

    return this.fetch(uri)
      .then(json => (validBoids(json) ? json.boids : []))
      .then(boids => boids[0]);
  }

  getInventory(boid) {
    const uri = this.buildQuery("inventory", { boid: boid });

    return this.fetch(uri)
      .then(json => json.inventory)
      .then(inventory => {
        if (!inventory) {
          return undefined;
        }

        const grouped = inventory.reduce((acc, item) => {
          const id = item.alt_link === 0 ? item.boid : item.alt_link;
          const group = (acc[id] = acc[id] || { boids: [] });

          group.quantity = Number(item.quantity);
          group.boids.push(item.boid);

          return acc;
        }, {});

        return Object.keys(grouped).map(key => grouped[key]);
      });
  }

  getModelInfo(boid) {
    const uri = this.buildQuery("lookup", { boid: boid });

    return this.fetch(uri).then(json => {
      if (!json || json.error) {
        return undefined;
      }

      const setNumbers = json.ids
        .filter(id => id.type === "set_number")
        .map(id => id.id);

      return {
        boid: boid,
        name: json.name,
        url: json.url,
        setNumbers: [...new Set(setNumbers)]
      };
    });
  }
}

export default Owl;
