class Inventory {
  constructor(api, setStorage, notifier) {
    const publishCompleted = (setNumber, inventory) => {
      const event = {
        eventType: "MODEL_INVENTORY_COMPLETE",
        setNumber: setNumber,
        inventory: inventory
      };
      console.log(`Publishing ${event.eventType} for set ${setNumber}`);
      return notifier.publish(event);
    };

    this.updateInventory = (setNumber, force) => {
      return setStorage.read(setNumber).then(storeInventory => {
        if (storeInventory && !force) {
          console.log(
            `Set ${setNumber} found in storage, and force was false.`
          );
          return Promise.resolve();
        }

        return api.getInventory(setNumber).then(inventory => {
          if (!inventory || inventory.length === 0) {
            console.log(`Got no inventory for set ${setNumber}`);
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
