import { Lambda } from "aws-sdk";

const publishMessage = (client, lambdaName, message) => {
  if (!message.eventType || message.eventType === "") {
    throw new Error("Missing required 'eventType' property");
  }

  const request = {
    FunctionName: lambdaName,
    Payload: JSON.stringify({
      body: message
    })
  };

  return client.invoke(request).promise();
};

class Notifier {
  constructor({ lambdaName, client = new Lambda() }) {
    this.publish = message => publishMessage(client, lambdaName, message);
  }
}

export default Notifier;
