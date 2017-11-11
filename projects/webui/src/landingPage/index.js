import React from 'react'
import { Col } from 'react-bootstrap'
import { Route } from 'react-router'
import Sidebar from '../headerBar'
import ModelPage from '../modelPage'

const LandingPage = ({ match }) =>
  <div>
    <Sidebar match={match} />
    <Route path="/model/:name" component={ModelPage} />
  </div>

export default LandingPage
