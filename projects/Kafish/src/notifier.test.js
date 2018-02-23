import Notifier from "./notifier";

const topic = "someTopic";
let publish, notifier, message;

beforeEach(() => {
  publish = jest.fn(dto => {
    message = dto;
    return { promise: () => Promise.resolve() };
  });

  notifier = new Notifier({
    topic: topic,
    client: { publish: publish }
  });
});

it("should send to the right topic", () =>
  notifier
    .publish({ eventType: "TEST_EVENT" })
    .then(() => expect(message.TopicArn).toEqual(topic)));

it("should throw if there is no eventType", () =>
  expect(() => notifier.publish({ wat: "no" })).toThrow());

it("should throw if there is a blank eventType", () =>
  expect(() => notifier.publish({ eventType: "", wat: "no" })).toThrow());

it("should throw if there is a non-string eventType", () =>
  expect(() => notifier.publish({ eventType: 123, wat: "no" })).toThrow());

it("should remove the eventType from message body", () =>
  notifier
    .publish({ eventType: "TEST", wat: "yes" })
    .then(() => expect(JSON.parse(message.Message)).toEqual({ wat: "yes" })));

it("should add eventType as a header", () =>
  notifier.publish({ eventType: "TEST", wat: "yes" }).then(() =>
    expect(message.MessageAttributes).toEqual({
      EventType: { DataType: "String", StringValue: "TEST" }
    })
  ));
