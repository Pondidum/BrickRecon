import React from 'react'
import { Table, Thead, Th } from 'reactable'

const imageSize = {
  width: 66,
  height: 50
}
const columns = [
  { key: 'image', label: 'Image', ...imageSize },
  { key: 'partNumber', label: 'Part No' },
  { key: 'name', label: 'Name' },
  { key: 'color', label: 'Color' },
  { key: 'quantity', label: 'Quantity' },
  { key: 'category', label: 'Category' }
]

const header = (col, i) =>
  <Th key={i} column={col.key} width={col.width ? col.width : null}>
    <strong className="name-header">
      {col.label}
    </strong>
  </Th>

const PartList = ({ parts }) => {
  const withImages = parts.map(part =>
    Object.assign({}, part, {
      image: (
        <img
          src={`/images/parts/${part.partNumber}-${part.color}.png`}
          alt={`part ${part.partNumber}`}
          width={imageSize.width + 'px'}
          height={imageSize.height + 'px'}
        />
      )
    })
  )
  return (
    <Table className="table table-hover" sortable={true} data={withImages}>
      <Thead>
        {columns.map(header)}
      </Thead>
    </Table>
  )
}

export default PartList
