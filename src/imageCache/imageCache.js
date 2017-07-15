const s3storage = require('./s3storage')
const fetchImage = require('./fetchImage')

const buildKey = part =>
  `website\\images\\parts\\${part.partno}-${part.colour}.png`

const buildUrl = part =>
  `https://img.bricklink.com/ItemImage/PN/${part.colour}/${part.partno}.png`

const put = (part, fetch = fetchImage, storage = s3storage) => {
  const key = buildKey(part)

  return storage
    .exists(key)
    .then(exists => {
      if (!exists) return fetch(buildUrl(part))
    })
    .then(image => {
      if (image) return storage.write(key, image)
    })
    .catch(err => console.log(err))
}

module.exports = {
  buildUrl: buildUrl,
  buildKey: buildKey,
  put: put
}
