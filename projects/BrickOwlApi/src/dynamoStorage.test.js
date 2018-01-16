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
