import { DynamoDB } from "aws-sdk";

const writeItem = (client, tableName, item) => {
  const request = {
    TableName: tableName,
    Item: item
  };

  return client.put(request).promise();
};

const readItem = (client, tableName, setNumber) => {
  const request = {
    TableName: tableName,
    Key: { setNumber: setNumber }
  };

  return client
    .get(request)
    .promise()
    .then(record => record.Item);
};

class SetStorage {
  constructor({ tableName, client = new DynamoDB.DocumentClient() }) {
    this.read = setNumber => readItem(client, tableName, setNumber);
    this.write = item => writeItem(client, tableName, item);
  }
}

export default SetStorage;
