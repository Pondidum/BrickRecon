import React from 'react'
import { Link, Route } from 'react-router-dom'
import path from 'path'

const ModelLink = ({ modelPath }) => {
  const name = path.basename(modelPath, path.extname(modelPath))
  const link = '/model/' + path.basename(modelPath)

  return (
    <Route
      path={link}
      exact={true}
      children={({ match }) =>
        <li className={match ? 'active' : ''}>
          <Link to={link}>
            {name}
          </Link>
        </li>}
    />
  )
}

export default ModelLink
