class Inventory {
  constructor(store, owl, notifier) {
    this.store = store;
    this.owl = owl;
    this.notifier = notifier;
  }

  buildModel(setNumber) {
    const owl = this.owl;
    return owl.getBoid(setNumber).then(setBoid => {
      if (!setBoid) {
        return undefined;
      }

      return owl
        .getModelInfo(setBoid)
        .then(model =>
          owl
            .getInventory(setBoid)
            .then(inv =>
              Object.assign({ setNumber: setNumber }, model, { inventory: inv })
            )
        );
    });
  }

  updateInventory(setNumber, force) {
    return this.store.read(setNumber).then(storeModel => {
      if (storeModel && !force) {
        return Promise.resolve();
      }

      return this.buildModel(setNumber).then(model => {
        if (!model) {
          return Promise.resolve();
        }

        return this.store.write(model).then(() =>
          this.notifier.publish({
            eventType: "MODEL_INVENTORY_COMPLETE",
            model
          })
        );
      });
    });
  }
}

export default Inventory;
