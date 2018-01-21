import { DynamoDB } from "aws-sdk";
import { chunk, mapFrom } from "./util";

const getMany = (client, tableName, batchSize, boids) => {
  const createRequest = group => ({
    RequestItems: {
      [tableName]: { Keys: group.map(boid => ({ boid: boid })) }
    }
  });

  const requests = chunk(boids, batchSize).map(group =>
    client
      .batchGet(createRequest(group))
      .promise()
      .then(data => data.Responses[tableName])
      .then(results => mapFrom(results, x => x.boid, x => x.partNumber))
  );

  return Promise.all(requests).then(results => Object.assign({}, ...results));
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

class DynamoStorage {
  constructor(tableName, options) {
    const client = options.client || new DynamoDB.DocumentClient();
    const batchSize = options.batchSize || 100;

    this.getMany = boids => getMany(client, tableName, batchSize, boids);
    this.writeMany = boids => writeMany(client, tableName, batchSize, boids);
  }
}

export default DynamoStorage;