const fetch = require("node-fetch");

const brickowl = `https://api.brickowl.com/v1`;
const key = process.env.BRICKOWL_TOKEN;

export const getBoid = setId => {
  const uri = `${brickowl}/catalog/id_lookup?key=${key}&type=Set&id=${setId}`;

  return fetch(uri)
    .then(res => res.json())
    .then(json => (json.boids.length > 0 ? Number(json.boids[0]) : undefined));
};
