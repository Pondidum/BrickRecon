const imageCache = require('./imageCache')

exports.cache = imageCache
exports.handler = (event, context, callback) => {
  const cache = exports.cache

  const cacheOperations = event.parts
    .map(cache.put)
    .map(p => p.catch(err => err))

  return Promise.all(cacheOperations)
    .then(results => callback(null, null))
    .catch(err => {
      console.error(err)
      callback(err)
    })
}
