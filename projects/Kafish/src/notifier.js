import { SNS, config } from "aws-sdk";

const stripEventType = contents => {
  const { eventType, ...model } = contents;
  return model;
};

const publish = (client, topic, contents) => {
  if (!contents.eventType || typeof contents.eventType !== "string") {
    throw new Error("You must specifiy an eventType string property");
  }

  const request = {
    TopicArn: topic,
    Message: JSON.stringify(stripEventType(contents)),
    MessageAttributes: {
      EventType: { DataType: "String", StringValue: contents.eventType }
    }
  };

  return client.publish(request).promise();
};

export default class Notification {
  constructor({ topic, client = new SNS() } = {}) {
    this.publish = contents => publish(client, topic, contents);
  }
}
