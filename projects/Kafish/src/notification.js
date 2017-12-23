import { SNS, config } from "aws-sdk";
const sns = new SNS();

const publish = (topic, contents) =>
  sns.publish({ TopicArn: topic, Message: message }).promise();

module.exports = topic => {
  return {
    publish: contents => publish(topic, contents)
  };
};
