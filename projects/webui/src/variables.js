import variables from './variables.dev.json'

export default {
    bucket: variables.bucket,
    environment: variables.environment,
    s3url: `http://${variables.bucket}.s3-eu-west-1.amazonaws.com/`
}
