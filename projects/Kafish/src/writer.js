import enhance from "./enhance";
import log from "./log";

export default options =>
  options.store
    .write(enhance(options.awsEvent.body))
    .then(() => options.respond("200", {}))
    .catch(err => {
      log.error("error writing to dynamo", err);

      options.respond("400", {
        message: "Unable to store event",
        exception: err
      });
    });
