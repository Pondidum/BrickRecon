import React from 'react'
import { Navbar } from 'react-bootstrap'
import { Link } from 'react-router-dom'
import ModelPicker from '../modelPicker'

const sidebar = ({ match }) =>
  <div className="row">
    <Navbar>
      <Navbar.Header>
        <Navbar.Brand>
          <Link to="/">BrickRecon</Link>
        </Navbar.Brand>
      </Navbar.Header>
      <ModelPicker match={match} />
    </Navbar>
  </div>

export default sidebar
