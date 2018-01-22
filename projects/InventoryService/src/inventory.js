class Inventory {
  constructor(api, setStorage, notifier) {
    const publishCompleted = (setNumber, inventory) =>
      notifier.publish({
        eventType: "MODEL_INVENTORY_COMPLETE",
        setNumber: setNumber,
        inventory: inventory
      });

    this.updateInventory = (setNumber, force) => {
      return setStorage.read(setNumber).then(storeInventory => {
        if (storeInventory && !force) {
          return Promise.resolve();
        }

        return api.getInventory(setNumber).then(inventory => {
          if (!inventory || inventory.length === 0) {
            return Promise.resolve();
          }
          return setStorage
            .write({ setNumber: setNumber, inventory: inventory })
            .then(() => publishCompleted(setNumber, inventory));
        });
      });
    };
  }
}

export default Inventory;
