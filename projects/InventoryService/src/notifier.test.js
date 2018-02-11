import Notifier from "./notifier";

const lambdaName = "wat_lambda";
let notification, notifier;

beforeEach(() => {
  const client = {
    invoke: m => {
      notification = m;
      return { promise: () => Promise.resolve() };
    }
  };

  notifier = new Notifier({
    lambdaName: lambdaName,
    client: client
  });
});

it("should publish a well formed message", () => {
  const message = { eventType: "TEST_MESSAGE", setNumber: 123 };

  return notifier.publish(message).then(() =>
    expect(notification).toEqual({
      FunctionName: lambdaName,
      Payload: '{"body":{"eventType":"TEST_MESSAGE","setNumber":123}}'
    })
  );
});

it("should reject messages without an eventType", () => {
  const message = { setNumber: 456 };

  expect(() => notifier.publish(message)).toThrow(
    "Missing required 'eventType' property"
  );
});
