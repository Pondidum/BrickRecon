import { DynamoDB } from "aws-sdk";

const getMany = boids => Promise.resolve({});
const writeMany = (client, tableName, batchSize, boids) => {
  const request = {
    RequestItems: {
      [tableName]: Object.keys(boids).map(boid => ({
        PutRequest: {
          Item: {
            boid: boid,
            partNumber: boids[boid]
          }
        }
      }))
    }
  };

  return client.batchWrite(request).promise();
};

class Storage {
  constructor(tableName, options) {
    const client = options.client || new DynamoDB.DocumentClient();
    const batchSize = options.batchSize || 100;

    this.getMany = boids => getMany(client, tableName, boids);
    this.writeMany = boids => writeMany(client, tableName, batchSize, boids);
  }
}

export default Storage;
