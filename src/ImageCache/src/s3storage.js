const aws = require('aws-sdk')
const s3 = new aws.S3()

const exists = (bucket, key) => {
  const query = {
    Bucket: bucket,
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

const write = (bucket, key, contents) => {
  const command = {
    Bucket: bucket,
    Key: key,
    Body: contents,
    ContentType: 'image/png'
  }

  return s3.putObject(command).promise()
}

module.exports = bucket => {
  return {
    exists: key => exists(bucket, key),
    write: (key, contents) => write(bucket, key, contents)
  }
}
