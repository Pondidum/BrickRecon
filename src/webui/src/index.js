import React from 'react'
import { render } from 'react-dom'

import { createStore, applyMiddleware } from 'redux'
import createHistory from 'history/createHashHistory'

import { Provider } from 'react-redux'
import { Route } from 'react-router'
import { ConnectedRouter, routerMiddleware } from 'react-router-redux'
import { apiMiddleware } from 'redux-api-middleware'

import reducers from './reducers'
import initialise from './initialise'
import registerServiceWorker from './registerServiceWorker'

import LandingPage from './landingPage'

const history = createHistory()
const middleware = routerMiddleware(history)

const devTools =
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
const createStoreWithMiddleware = applyMiddleware(apiMiddleware, middleware)(
  createStore
)
const store = createStoreWithMiddleware(reducers, devTools)

initialise(store.dispatch)

render(
  <Provider store={store}>
    <ConnectedRouter history={history}>
      <div>
        <Route path="/" exact component={LandingPage} />
      </div>
    </ConnectedRouter>
  </Provider>,
  document.getElementById('root')
)

registerServiceWorker()
