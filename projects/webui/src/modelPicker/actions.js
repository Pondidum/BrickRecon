import { CALL_API } from 'redux-api-middleware'
import variables from '../variables'

export const listModels = () => {
  return {
    [CALL_API]: {
      endpoint: variables.s3url + '?prefix=models/',
      method: 'GET',
      types: [
        'LIST_ALL_MODELS_REQUEST',
        {
          type: 'LIST_ALL_MODELS_SUCCESS',
          payload: (action, state, res) => res.text()
        },
        'LIST_ALL_MODELS_FAILURE'
      ]
    }
  }
}

export const loadModel = modelName => {
  return {
    [CALL_API]: {
      endpoint: variables.s3url + 'models/' + modelName + '.json',
      method: 'GET',
      types: ['LOAD_MODEL_REQUEST', 'LOAD_MODEL_SUCCESS', 'LOAD_MODEL_FAILURE']
    }
  }
}
