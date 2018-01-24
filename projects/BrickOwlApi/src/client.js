import fetch from "node-fetch";
import queryString from "query-string";
import { chunk, mapFrom } from "./util";

const buildQuery = (path, token, query) => {
  const map = Object.assign({ key: token }, query);
  const qs = queryString.stringify(map);

  return `https://api.brickowl.com/v1/catalog/${path}?${qs}`;
};

const combineAlternateParts = inventory => {
  const grouped = inventory.reduce((acc, item) => {
    const id = item.alt_link === 0 ? item.boid : item.alt_link;
    const group = (acc[id] = acc[id] || { boids: [] });

    group.quantity = Number(item.quantity);
    group.boids.push(item.boid);

    return acc;
  }, {});

  return Object.keys(grouped).map(key => grouped[key]);
};

const bestPartNumber = ids => {
  const map = ids.reduce((all, current) => {
    all[current.type] = all[current.type] || [];
    all[current.type].push(current.id);
    return all;
  }, {});

  if (map.ldraw) return map.ldraw[0];
  if (map.peeron_id) return map.peeron_id[0];
  if (map.design_id) return map.design_id[0];
};

//---------------------------------------------------------------------------//

const getSetBoid = (fetcher, token, setNumber) => {
  const uri = buildQuery("id_lookup", token, { type: "Set", id: setNumber });
  const validBoids = js => js && !js.error && js.boids && js.boids.length > 0;

  return fetcher(uri)
    .then(json => (validBoids(json) ? json.boids : []))
    .then(boids => boids[0]);
};

const getInventory = (fetcher, token, boid) => {
  const uri = buildQuery("inventory", token, { boid: boid });

  return fetcher(uri)
    .then(json => json.inventory)
    .then(inventory => (inventory ? combineAlternateParts(inventory) : []));
};

const getPartNumbers = (fetcher, batchSize, token, boids) => {
  const chunks = chunk(
    boids.map(boid => boid.replace(/(-.*)/g, "")),
    batchSize
  );

  const queries = chunks.map(groups => {
    const uri = buildQuery("bulk_lookup", token, { boids: groups.join(",") });
    return fetcher(uri)
      .then(response => response.items)
      .then(parts =>
        mapFrom(Object.keys(parts), x => x, x => bestPartNumber(parts[x].ids))
      );
  });

  return Promise.all(queries).then(maps => Object.assign({}, ...maps));
};

const defaultFetcher = uri => fetch(uri).then(res => res.json());

class Client {
  constructor(token, { fetcher = defaultFetcher, batchSize = 100 } = {}) {
    this.getSetBoid = setNumber => getSetBoid(fetcher, token, setNumber);
    this.getInventory = boid => getInventory(fetcher, token, boid);
    this.getPartNumbers = boids =>
      getPartNumbers(fetcher, batchSize, token, boids);
  }
}

export default Client;
