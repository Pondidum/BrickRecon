import { DynamoDB } from "aws-sdk";
import { chunk, mapFrom } from "./util";

const getMany = (client, tableName, boids) => {
  const createRequest = group => ({
    RequestItems: {
      [tableName]: { Keys: group.map(boid => ({ boid: boid })) }
    }
  });

  const uniqueBoids = [...new Set(boids).keys()];
  const requests = chunk(uniqueBoids, 100).map(group =>
    client
      .batchGet(createRequest(group))
      .promise()
      .then(data => data.Responses[tableName])
      .then(results => mapFrom(results, x => x.boid, x => x.partNumber))
  );

  return Promise.all(requests).then(results => Object.assign({}, ...results));
};

const writeMany = (client, tableName, boids) => {
  const requests = Object.keys(boids).map(boid => ({
    PutRequest: {
      Item: {
        boid: boid,
        partNumber: boids[boid]
      }
    }
  }));

  const chunks = chunk(requests, 25);

  return Promise.all(
    chunks.map(group =>
      client.batchWrite({ RequestItems: { [tableName]: group } }).promise()
    )
  );
};

class DynamoStorage {
  constructor(tableName, { client = new DynamoDB.DocumentClient() } = {}) {
    this.getMany = boids => getMany(client, tableName, boids);
    this.writeMany = boids => writeMany(client, tableName, boids);
  }
}

export default DynamoStorage;
