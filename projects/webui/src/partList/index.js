import React from "react";
import { Table, Thead, Th } from "reactable";
import naturalSort from "javascript-natural-sort";
import PartImage from "./partImage";
import PartLinks from "./partLinks";

const imageSize = {
  width: 66,
  height: 50
};

const buildHeader = col => (
  <Th key={col.key} column={col.key} width={col.width ? col.width : null}>
    <strong className="name-header">{col.label}</strong>
  </Th>
);

const headers = [
  { key: "image", label: "Image", ...imageSize },
  { key: "partNumber", label: "Part No" },
  { key: "name", label: "Name" },
  { key: "colorName", label: "Color" },
  { key: "quantity", label: "Quantity" },
  { key: "category", label: "Category" },
  { key: "links", label: "Links" }
].map(buildHeader);

const sort = [
  { column: "partNumber", sortFunction: naturalSort },
  { column: "name", sortFunction: naturalSort },
  { column: "colorName" },
  { column: "quantity", sortFunction: naturalSort },
  { column: "category" }
];

const defaultSort = {
  column: "name",
  direction: "asc"
};

const enrich = part =>
  Object.assign({}, part, {
    image: <PartImage part={part} />,
    links: <PartLinks part={part} />
  });

const PartList = ({ parts }) => (
  <Table
    className="table table-hover"
    data={parts.map(enrich)}
    sortable={sort}
    defaultSort={defaultSort}
  >
    <Thead>{headers}</Thead>
  </Table>
);

export default PartList;
