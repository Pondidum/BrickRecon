import reduce from './reducers'

const xml = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Name>brickrecon-dev</Name>
  <Prefix>models/</Prefix>
  <Marker></Marker>
  <MaxKeys>1000</MaxKeys>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>models/dual-rail-gun.json</Key>
    <LastModified>2017-07-19T19:40:49.000Z</LastModified>
    <ETag>&quot;1f199a369665791d8024ad0136cc1926&quot;</ETag>
    <Size>1233</Size>
    <StorageClass>STANDARD</StorageClass>
  </Contents>
</ListBucketResult>
  `
describe('LIST_ALL_MODELS_SUCCESS', () => {
  it('should handle bucket xml', () => {
    const event = {
      type: 'LIST_ALL_MODELS_SUCCESS',
      payload: xml
    }
    const state = reduce(undefined, event)

    expect(state).toEqual({
      selected: null,
      available: ['models/dual-rail-gun.json']
    })
  })
})

describe('LOAD_MODEL_REQUEST', () => {
  it('should clear out the existing model', () => {
    const event = { type: 'LOAD_MODEL_REQUEST' }
    const state = reduce({ available: [], selected: { wat: 'is this' } }, event)

    expect(state).toEqual({
      available: [],
      selected: null
    })
  })
})

describe('LOAD_MODEL_SUCCESS', () => {
  it('should update the state with the model', () => {
    const event = { type: 'LOAD_MODEL_SUCCESS', payload: { parts: [] } }
    const state = reduce(undefined, event)

    expect(state).toEqual({
      available: [],
      selected: { parts: [] }
    })
  })

  it('should hydrate in the colorName', () => {
    const event = {
      type: 'LOAD_MODEL_SUCCESS',
      payload: {
        parts: [
          {
            partNumber: '3039',
            name: 'Slope 45 2 x 2',
            color: 11,
            quantity: 0,
            category: 'Slope'
          },
          {
            partNumber: '14769',
            name: 'Tile, Round 2 x 2 with Bottom Stud Holder',
            color: 86,
            quantity: 2,
            category: 'Tile, Round'
          }
        ]
      }
    }
    const state = reduce(undefined, event)

    expect(state.selected.parts[0].colorName).toEqual('Black')
    expect(state.selected.parts[1].colorName).toEqual('Light Bluish Gray')
  })
})
