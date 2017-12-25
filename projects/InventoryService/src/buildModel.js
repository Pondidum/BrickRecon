const buildModel = (owl, setNumber) =>
  owl
    .getBoid(setNumber)
    .then(boid => Promise.all([owl.getModelInfo(boid), owl.getInventory(boid)]))
    .then(([model, inventory]) =>
      Object.assign({}, model, { inventory: inventory })
    );

export default buildModel;
