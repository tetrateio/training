const proxy = require('http-proxy-middleware');

// create-react-app v2 forces us to make this file
// source: https://github.com/facebook/create-react-app/blob/master/packages/react-scripts/template/README.md#configuring-the-proxy-manually
module.exports = function(app) {
  const listEndpoints = {
    '/v1': {
      // We setup dynamic proxy here so we can edit with environment variables.
      target:
        process.env.REACT_APP_PROXY || 'http://35.197.239.230',
      pathRewrite: {
        // We can do `pathRewrite` here if required. e.g. '/v1': '/v2'.
      },
      changeOrigin: true
    }
  };

  Object.keys(listEndpoints).forEach(endpoint => {
    app.use(proxy(endpoint, listEndpoints[endpoint]));
  });
};
