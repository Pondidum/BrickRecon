import React from "react";

const brickOwl = part =>
  `http://www.brickowl.com/search/catalog?query=${part.partNumber}&cat=1`;
const peeron = part => `http://peeron.com/inv/parts/${part.partNumber}`;

const PartLinks = ({ part }) => (
  <ul className="list-unstyled">
    <li>
      <a href={brickOwl(part)}>BrickOwl</a>
    </li>
    <li>
      <a href={peeron(part)}>Peeron</a>
    </li>
  </ul>
);

export default PartLinks;
