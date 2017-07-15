const aws = require('aws-sdk')
const s3 = new aws.S3()

const exists = key => {
  const query = {
    Bucket: 'brickrecon',
    Key: key
  }

  return new Promise((resolve, reject) =>
    s3
      .headObject(query)
      .promise()
      .then(data => resolve(true))
      .catch(err => resolve(false))
  )
}

const write = (key, contents) => {
  const command = {
    Bucket: 'brickrecon',
    Key: key,
    Body: contents,
    ContentType: 'image/png'
  }

  return s3.putObject(command).promise()
}

module.exports = {
  exists: key => exists,
  write: write
}
