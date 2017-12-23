import enhance from "./enhance";
import { error } from "./log";

export default options => {
  const event = enhance(options.awsEvent.body);
  return options.store
    .write(event)
    .then(() => options.publish(event))
    .then(() => options.respond("200", {}))
    .catch(err => {
      error("error writing to dynamo", err);

      options.respond("400", {
        message: "Unable to store event",
        exception: err
      });
    });
};
