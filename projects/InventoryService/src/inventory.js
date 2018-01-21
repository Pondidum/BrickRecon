class Inventory {
  constructor(api, setStorage, notifier) {
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
            .write({
              setNumber: setNumber,
              inventory: inventory
            })
            .then(() =>
              notifier.publish({
                eventType: "MODEL_INVENTORY_COMPLETE",
                setNumber: setNumber,
                inventory: inventory
              })
            );
        });
      });
    };
  }
}

export default Inventory;
