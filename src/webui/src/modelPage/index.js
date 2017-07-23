import React from 'react'
import { Col } from 'react-bootstrap'
import { connect } from 'react-redux'
import { loadModel } from '../modelPicker/actions'

const mapStateToProps = (state, ownProps) => {
  return {
    ...ownProps,
    model: state.models.selected
  }
}

const mapDispatchToProps = dispatch => {
  return {
    loadModel: name => dispatch(loadModel(name))
  }
}

const ModelPage = ({ model, loadModel, match }) => {
  const modelName = match.params.name
  if (!model || model.name !== modelName) {
    loadModel(modelName)
  }

  return (
    <div className="row">
      <Col sm={12} className="main">
        <h1>
          {model ? model.name : 'none'}
        </h1>
      </Col>
    </div>
  )
}

export default connect(mapStateToProps, mapDispatchToProps)(ModelPage)
