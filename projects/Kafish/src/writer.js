import enhance from "./enhance";

const write = (dynamo, tableName, event) =>
  dynamo.put({ TableName: tableName, Item: event }).promise();

const error = (message, err) =>
  console.log(message, JSON.stringify(err, null, 2));

export default options =>
  write(options.dynamo, options.tableName, enhance(options.awsEvent.body))
    .then(() => options.respond("200", {}))
    .catch(err => {
      error("error writing to dynamo", err);

      options.respond("400", {
        message: "Unable to store event",
        exception: err
      });
    });
