const lambda = require('./index')

it('should handle multiple part lookups', () => {
  const event = {
    parts: [{ partno: 1, color: 'a' }, { partno: 2, color: 'b' }]
  }

  const callback = jest.fn()

  lambda.cache = {
    put: jest.fn(() => Promise.resolve())
  }

  return lambda.handler(event, {}, callback).then(() => {
    expect(callback).toHaveBeenCalledWith(null, null)
  })
})
