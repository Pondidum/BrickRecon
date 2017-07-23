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
    const event = { type: 'LOAD_MODEL_SUCCESS', payload: { some: 'thing' } }
    const state = reduce(undefined, event)

    expect(state).toEqual({
      available: [],
      selected: { some: 'thing' }
    })
  })
})
