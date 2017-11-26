import React from 'react'

const defaultImageSize = {
  width: 66,
  height: 50
}

const PartImage = ({ part, imageSize = defaultImageSize }) => <img
  src={part.imageUrl}
  alt={`part ${part.partNumber}`}
  width={imageSize.width + 'px'}
  height={imageSize.height + 'px'}
/>

export default PartImage
