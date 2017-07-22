import { CALL_API } from 'redux-api-middleware'

export const listModels = () => {
  return {
    [CALL_API]: {
      endpoint:
        'http://brickrecon-dev.s3-eu-west-1.amazonaws.com/?prefix=models/',
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
