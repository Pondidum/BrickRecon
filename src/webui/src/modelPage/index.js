import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Col } from 'react-bootstrap'
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

class ModelPage extends Component {
  constructor(params) {
    super(params)
    const { model, loadModel, match } = params

    const modelName = match.params.name

    this.modelName = modelName
    this.model = model
    this.loadModel = () => loadModel(modelName)
  }

  componentWillMount = () => {
    if (this.model && this.model.name === this.modelName) return

    this.loadModel()
  }

  render() {
    return (
      <div className="row">
        <Col sm={12} className="main">
          <h1>
            {this.model ? this.model.name : 'none'}
          </h1>
        </Col>
      </div>
    )
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(ModelPage)
