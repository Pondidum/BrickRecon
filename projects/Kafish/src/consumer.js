import rebuild from "./rebuild";

export const consumerHandler = (awsEvent, context, callback) => {
  const publisher = {};
  awsEvent.Records.map(record => record.dynamodb.NewImage)
    .map(element => rebuild(element))
    .forEach(record => publisher.publish(record));
};
