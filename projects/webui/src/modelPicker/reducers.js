import xpath from 'xpath'
import { DOMParser as dom } from 'xmldom'
import colors from '../domain/colors'

const defaultState = {
  available: [],
  selected: null
}

const transformModelList = xml => {
  const doc = new dom().parseFromString(xml)
  const ns = 'http://s3.amazonaws.com/doc/2006-03-01/'
  const path = '/ns:ListBucketResult/ns:Contents/ns:Key/text()'

  const select = xpath.useNamespaces({ ns: ns })
  const result = select(path, doc).map(x => x.nodeValue)

  return result
}

const hydrateModel = model => {
  const partsWithColor = model.parts.map(part =>
    Object.assign({ colorName: colors[part.color] }, part)
  )

  return Object.assign({}, model, { parts: partsWithColor })
}

export default (state = defaultState, action) => {
  switch (action.type) {
    case 'LIST_ALL_MODELS_REQUEST': {
      return state
    }

    case 'LIST_ALL_MODELS_SUCCESS': {
      return Object.assign({}, state, {
        available: transformModelList(action.payload)
      })
    }

    case 'LIST_ALL_MODELS_FAILURE': {
      return state
    }

    case 'LOAD_MODEL_REQUEST': {
      return Object.assign({}, state, { selected: null })
    }

    case 'LOAD_MODEL_SUCCESS': {
      return Object.assign({}, state, {
        selected: hydrateModel(action.payload)
      })
    }

    case 'LOAD_MODEL_FAILURE': {
      return Object.assign({}, state, { selected: null })
    }

    default: {
      return state
    }
  }
}
