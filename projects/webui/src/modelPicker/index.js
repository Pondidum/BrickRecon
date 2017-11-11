import React from 'react'
import { Nav } from 'react-bootstrap'
import { connect } from 'react-redux'
import ModelLink from './modelLink'

const mapStateToProps = (state, ownProps) => {
  return {
    models: state.models.available
  }
}

const picker = ({ models, onSelect = () => {} }) =>
  <Nav activeKey="1" onSelect={onSelect}>
    {models.map((modelPath, i) => <ModelLink key={i} modelPath={modelPath} />)}
  </Nav>

export default connect(mapStateToProps)(picker)
