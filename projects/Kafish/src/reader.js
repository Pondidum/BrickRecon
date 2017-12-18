const error = (message, err) =>
  console.log(message, JSON.stringify(err, null, 2));

export default options =>
  options.store
    .read(options.awsEvent.queryStringParameters.timestamp)
    .then(data => options.respond("200", data.Items))
    .catch(err => {
      error("error reading from dynamo", err);

      options.respond("400", {
        message: "Unable to read events",
        exception: err
      });
    });
