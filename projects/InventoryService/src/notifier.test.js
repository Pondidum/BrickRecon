import Notifier from "./notifier";

it("should publish a well formed message", () => {
  let message;
  const client = {
    invoke: m => {
      message = m;
      return { promise: () => Promise.resolve() };
    }
  };

  const notifier = new Notifier({
    lambdaName: "wat",
    client: client
  });

  return notifier.publish({ setNumber: 123 }).then(() =>
    expect(message).toEqual({
      FunctionName: "wat",
      Payload: '{"setNumber":123}'
    })
  );
});
