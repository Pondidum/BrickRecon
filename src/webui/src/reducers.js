import { routerReducer } from 'react-router-redux'
import { combineReducers } from 'redux'
import models from './modelPicker/reducers'

export default combineReducers({
  router: routerReducer,
  models: models
})
