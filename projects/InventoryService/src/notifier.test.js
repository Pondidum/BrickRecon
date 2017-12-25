import Notifier from "./notifier";

it("should publish a well formed message", () => {
  let message;
  const client = {
    publish: m => {
      message = m;
      return { promise: () => Promise.resolve() };
    }
  };

  const notifier = new Notifier("wat", client);

  return notifier.publish({ setNumber: 123 }).then(() =>
    expect(message).toEqual({
      TopicArn: "wat",
      Message: '{"setNumber":123}'
    })
  );
});
