const fetchImage = require('./fetchImage')

it('should fetch known images', () => {
  const url = 'https://img.bricklink.com/ItemImage/PN/1/3022.png'
  return fetchImage(url).then(data => expect(data).toBeDefined())
})

it('should throw on non-existing', () => {
  const url = 'https://img.bricklink.com/ItemImage/wefwf/1/ef32f232f3f3f.png'
  return fetchImage(url).then(data => expect(data).toBeFalsy())
})
