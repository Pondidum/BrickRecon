import xpath from 'xpath'
import { DOMParser as dom } from 'xmldom'

const defaultState = []

const transform = xml => {
  const doc = new dom().parseFromString(xml)
  const ns = 'http://s3.amazonaws.com/doc/2006-03-01/'
  const path = '/ns:ListBucketResult/ns:Contents/ns:Key/text()'

  const select = xpath.useNamespaces({ ns: ns })
  const result = select(path, doc).map(x => x.nodeValue)

  return result
}

export default (state = defaultState, action) => {
  switch (action.type) {
    case 'LIST_ALL_MODELS_REQUEST': {
      return state
    }

    case 'LIST_ALL_MODELS_SUCCESS': {
      return transform(action.payload)
    }

    case 'LIST_ALL_MODELS_FAILURE': {
      return state
    }

    default: {
      return state
    }
  }
}
