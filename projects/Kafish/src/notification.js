import { SNS, config } from "aws-sdk";
const sns = new SNS();

const publish = (topic, contents) =>
  sns.publish({ TopicArn: topic, Message: JSON.stringify(contents) }).promise();

module.exports = topic => {
  return {
    publish: contents => publish(topic, contents)
  };
};
