import { DynamoDB } from "aws-sdk";

class Storage {
  constructor(configuration) {
    this.tableName = configuration.tableName;
    this.hashKey = configuration.hashKey;
    this.dynamo = new DynamoDB.DocumentClient({
      endpoint: configuration.endpoint
    });
  }

  writeMany(items) {
    const request = {
      RequestItems: {
        [this.tableName]: items.map(item => {
          return {
            PutRequest: {
              Item: item
            }
          };
        })
      }
    };

    return this.dynamo.batchWrite(request).promise();
  }

  write(item) {
    const request = {
      TableName: this.tableName,
      Item: item
    };
    return this.dynamo.put(request).promise();
  }

  read(key) {
    const request = {
      TableName: this.tableName,
      Key: {
        [this.hashKey]: key
      }
    };

    return this.dynamo
      .get(request)
      .promise()
      .then(record => record.Item);
  }

  readMany(filter) {
    return this.dynamo
      .scan({
        TableName: this.tableName,
        FilterExpression: filter.expression,
        ExpressionAttributeNames: filter.nameMap,
        ExpressionAttributeValues: filter.valueMap
      })
      .promise()
      .then(records => records.Items);
  }
}

export default Storage;
