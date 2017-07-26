const fetch = require('node-fetch')

const forPart = partNumber => {
  const searchUrl = `http://www.brickowl.com/search/catalog?query=${partNumber}&cat=1`

  return fetch(searchUrl)
    .then(res => res.text())
    .then(body => console.log(body))
}

module.exports = {
  part: forPart
}
