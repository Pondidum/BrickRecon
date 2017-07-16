const imageCache = require('./imageCache')

const part = { partno: 123, colour: 1 }

it('should not fetch an image if one is in storage', () => {
  const storage = {
    exists: () => Promise.resolve(true),
    write: jest.fn(() => Promise.resolve())
  }

  const fetchImage = jest.fn()

  return imageCache.put(fetchImage, storage, part).then(() => {
    expect(storage.write).not.toHaveBeenCalled()
    expect(fetchImage).not.toHaveBeenCalled()
  })
})

it('should store an image if one doesnt exist in storage', () => {
  const storage = {
    exists: () => Promise.resolve(false),
    write: jest.fn(() => Promise.resolve())
  }

  const image = new Buffer(0)
  const fetchImage = jest.fn(() => image)

  return imageCache.put(fetchImage, storage, part).then(() => {
    expect(storage.write).toHaveBeenCalledWith(imageCache.buildKey(part), image)
  })
})

it('should not store an image if it fails to be fetched', () => {
  const storage = {
    exists: () => Promise.resolve(false),
    write: jest.fn(() => Promise.resolve())
  }

  const fetchImage = jest.fn(() => null)

  return imageCache.put(fetchImage, storage, part).then(() => {
    expect(storage.write).not.toHaveBeenCalled()
  })
})
