import { SNS } from "aws-sdk";

class Notifier {
  constructor(topic, client) {
    this.topic = topic;
    this.client = client || new SNS();
  }

  publish(message) {
    return this.client
      .publish({ TopicArn: this.topic, Message: JSON.stringify(message) })
      .promise();
  }
}

export default Notifier;
