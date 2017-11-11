const pageFor = require('./pageFor')

describe('pageFor.part', () => {
  it('should return the url for a known part', () => {
    pageFor
      .part(30374)
      .then(url =>
        expect(url).toEqual(
          'http://www.brickowl.com/catalog/lego-bar-4l-1-x-4-21462-30374'
        )
      )
  })
})
