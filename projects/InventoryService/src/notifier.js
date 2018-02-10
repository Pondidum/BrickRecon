import { Lambda } from "aws-sdk";

const publishMessage = (client, lambdaName, message) =>
  client
    .invoke({
      FunctionName: lambdaName,
      Payload: JSON.stringify(message)
    })
    .promise();

class Notifier {
  constructor({ lambdaName, client = new Lambda() }) {
    this.publish = message => publishMessage(client, lambdaName, message);
  }
}

export default Notifier;
