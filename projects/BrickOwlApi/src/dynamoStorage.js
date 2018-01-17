import { DynamoDB } from "aws-sdk";
import { chunk } from "./util";

const getMany = (client, tableName, batchSize, boids) => {
  const request = {
    RequestItems: {
      [tableName]: {
        Keys: boids.map(boid => ({ boid: boid }))
      }
    }
  };

  return client
    .batchGet(request)
    .promise()
    .then(data => data.Responses[tableName])
    .then(results =>
      results.reduce((all, current) => {
        all[current.boid] = current.partNumber;
        return all;
      }, {})
    );
};

const writeMany = (client, tableName, batchSize, boids) => {
  const requests = Object.keys(boids).map(boid => ({
    PutRequest: {
      Item: {
        boid: boid,
        partNumber: boids[boid]
      }
    }
  }));

  const chunks = chunk(requests, batchSize);

  return Promise.all(
    chunks.map(group =>
      client.batchWrite({ RequestItems: { [tableName]: group } }).promise()
    )
  );
};

class Storage {
  constructor(tableName, options) {
    const client = options.client || new DynamoDB.DocumentClient();
    const batchSize = options.batchSize || 100;

    this.getMany = boids => getMany(client, tableName, batchSize, boids);
    this.writeMany = boids => writeMany(client, tableName, batchSize, boids);
  }
}

export default Storage;
