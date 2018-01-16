import DynamoStorage from "./dynamoStorage";

let client, storage;

beforeEach(() => {
  client = {
    batchWrite: jest.fn()
  };

  client.batchWrite.mockReturnValue({
    promise: () => Promise.resolve()
  });

  storage = new DynamoStorage("wat", {
    client: client,
    batchSize: 5
  });
});

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
