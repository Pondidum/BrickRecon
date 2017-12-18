import enhance from "./enhance";

const error = (message, err) =>
  console.log(message, JSON.stringify(err, null, 2));

export default options =>
  options.store
    .write(enhance(options.awsEvent.body))
    .then(() => options.respond("200", {}))
    .catch(err => {
      error("error writing to dynamo", err);

      options.respond("400", {
        message: "Unable to store event",
        exception: err
      });
    });
