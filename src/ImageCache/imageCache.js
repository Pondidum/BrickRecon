const buildKey = part => `images/parts/${part.partno}-${part.color}.png`

const buildUrl = part =>
  `https://img.bricklink.com/ItemImage/PN/${part.color}/${part.partno}.png`

const put = (fetch, storage, part) => {
  const key = buildKey(part)

  return storage
    .exists(key)
    .then(exists => {
      if (!exists) return fetch(buildUrl(part))
    })
    .then(image => {
      if (image) return storage.write(key, image)
    })
}

module.exports = {
  buildUrl: buildUrl,
  buildKey: buildKey,
  put: put
}
