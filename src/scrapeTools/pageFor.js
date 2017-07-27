const fetch = require('node-fetch')
const cheerio = require('cheerio')

const root = 'http://www.brickowl.com'

const forPart = partNumber => {
  const searchUrl = `${root}/search/catalog?query=${partNumber}&cat=1`

  return fetch(searchUrl).then(res => res.text()).then(body => {
    const doc = cheerio.load(body)
    const fragment = doc('ul.category-grid li a').attr('href')

    return root + fragment
  })
}

module.exports = {
  part: forPart
}
