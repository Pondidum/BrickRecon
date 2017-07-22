import React from 'react'
import { Col } from 'react-bootstrap'
import Sidebar from '../headerBar'

const LandingPage = () =>
  <div>
    <Sidebar />
    <div className="row">
      <Col sm={12} className="main">
        <h1>Content</h1>
      </Col>
    </div>
  </div>

export default LandingPage
