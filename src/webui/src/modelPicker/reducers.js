import xpath from 'xpath'
import { DOMParser as dom } from 'xmldom'

const defaultState = {
  available: [],
  selected: null
}

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
      return Object.assign({}, state, { available: transform(action.payload) })
    }

    case 'LIST_ALL_MODELS_FAILURE': {
      return state
    }

    case 'LOAD_MODEL_REQUEST': {
      return Object.assign({}, state, { selected: null })
    }

    case 'LOAD_MODEL_SUCCESS': {
      return Object.assign({}, state, { selected: action.payload })
    }

    case 'LOAD_MODEL_FAILURE': {
      return Object.assign({}, state, { selected: null })
    }

    default: {
      return state
    }
  }
}
