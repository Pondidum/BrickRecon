import React from 'react'
import { Table, Thead, Th } from 'reactable'

const columns = [
  { key: 'partNumber', label: 'Part No' },
  { key: 'name', label: 'Name' },
  { key: 'color', label: 'Color' },
  { key: 'quantity', label: 'Quantity' },
  { key: 'category', label: 'Category' }
]

const header = (col, i) =>
  <Th key={i} column={col.key}>
    <strong className="name-header">
      {col.label}
    </strong>
  </Th>

const PartList = ({ parts }) =>
  <Table className="table table-hover" sortable={true} data={parts}>
    <Thead>
      {columns.map(header)}
    </Thead>
  </Table>

export default PartList
