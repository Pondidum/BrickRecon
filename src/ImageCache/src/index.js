const s3storage = require('./s3storage')
const fetchImage = require('./fetchImage')
const imageCache = require('./imageCache')

exports.cache = imageCache
exports.handler = (event, context, callback) => {
  const cache = exports.cache

  const cacheOperations = event.parts
    .map(part => cache.put(fetchImage, s3storage(process.env.bucket), part))
    .map(p => p.catch(err => console.error(err)))

  return Promise.all(cacheOperations).then(results => callback(null, null))
}
