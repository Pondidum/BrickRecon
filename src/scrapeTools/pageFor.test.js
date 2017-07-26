const pageFor = require('./pageFor')

describe('pageFor.part', () => {
  it('should return the url for a known part', () => {
    expect(pageFor.part(30374)).toEqual(
      'http://www.brickowl.com/catalog/lego-bar-4l-1-x-4-21462-30374'
    )
  })
})
