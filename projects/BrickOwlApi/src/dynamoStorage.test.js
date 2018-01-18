import DynamoStorage from "./dynamoStorage";
import AWS, { DynamoDB } from "aws-sdk";

let client, storage;

beforeEach(() => {
  client = {
    batchWrite: jest.fn(),
    batchGet: jest.fn()
  };

  client.batchWrite.mockReturnValue({
    promise: () => Promise.resolve()
  });

  storage = new DynamoStorage("wat", {
    client: client,
    batchSize: 5
  });
});

const dynamoReturns = result => {
  client.batchGet.mockReturnValueOnce({
    promise: () => Promise.resolve(result)
  });
};

describe("writeMany", () => {
  it("should call the client with the right structure", () =>
    storage.writeMany({ boid1: "part1", boid2: "part2" }).then(() =>
      expect(client.batchWrite.mock.calls[0][0]).toEqual({
        RequestItems: {
          wat: [
            { PutRequest: { Item: { boid: "boid1", partNumber: "part1" } } },
            { PutRequest: { Item: { boid: "boid2", partNumber: "part2" } } }
          ]
        }
      })
    ));

  it("should call the client multiple times with batches", () => {
    const seed = [...new Array(13).keys()];
    const request = seed.reduce((a, c) => {
      a["boid" + c] = "part" + c;
      return a;
    }, {});

    return storage.writeMany(request).then(() => {
      expect(client.batchWrite.mock.calls.length).toEqual(3);
      expect(client.batchWrite.mock.calls[0][0]).toEqual({
        RequestItems: {
          wat: seed.slice(0, 5).map(i => ({
            PutRequest: { Item: { boid: "boid" + i, partNumber: "part" + i } }
          }))
        }
      });
    });
  });
});

describe("getMany", () => {
  it("should return empty object when none found", () => {
    dynamoReturns({
      Responses: {
        wat: []
      }
    });

    return storage
      .getMany(["boid1", "boid4"])
      .then(data => expect(data).toEqual({}));
  });

  it("should only return boids which were found", () => {
    dynamoReturns({
      Responses: {
        wat: [{ boid: "boid1", partNumber: "part1" }]
      }
    });

    return storage.getMany(["boid1", "boid2"]).then(result =>
      expect(result).toEqual({
        boid1: "part1"
      })
    );
  });

  it("should batch the requests when over batchSize", () => {
    const seed = [...new Array(13).keys()];
    const range = (start, finish) =>
      seed
        .slice(start, finish)
        .map(i => ({ boid: "boid" + i, partNumber: "part" + i }));

    dynamoReturns({ Responses: { wat: range(0, 5) } });
    dynamoReturns({ Responses: { wat: range(5, 10) } });
    dynamoReturns({ Responses: { wat: range(10, 13) } });

    const request = seed.map(i => "boid" + i);
    const expected = seed.reduce((all, current) => {
      all["boid" + current] = "part" + current;
      return all;
    }, {});

    return storage.getMany(request).then(result => {
      expect(client.batchGet.mock.calls.length).toEqual(3);
      expect(result).toEqual(expected);
    });
  });
});
