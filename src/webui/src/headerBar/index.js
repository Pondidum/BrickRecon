import React from 'react'
import { Navbar } from 'react-bootstrap'
import { Link } from 'react-router-dom'
import ModelPicker from '../modelPicker'

const sidebar = () =>
  <div className="row">
    <Navbar>
      <Navbar.Header>
        <Navbar.Brand>
          <Link to="/">BrickRecon</Link>
        </Navbar.Brand>
      </Navbar.Header>
      <ModelPicker />
    </Navbar>
  </div>

export default sidebar
