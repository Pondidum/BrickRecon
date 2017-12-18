import log from "./log";

export default options =>
  options.store
    .read(options.awsEvent.queryStringParameters.timestamp)
    .then(data => options.respond("200", data.Items))
    .catch(err => {
      log.error("error reading from dynamo", err);

      options.respond("400", {
        message: "Unable to read events",
        exception: err
      });
    });
