const createProxyMiddleware = require('http-proxy-middleware');
const dev = !process.env.NODE_ENV || process.env.NODE_ENV === 'development';
const dockerhost = dev ? process.env.REACT_APP_DOCKER_PATH : 'localhost';
const port = 8080;
const url = `http://${dockerhost}:${port}`;

module.exports = function (app) {
  app.use(
    '/api/test/',
    createProxyMiddleware({
      target: `http://localhost:8000`,
      pathRewrite: { '^/api/test': '/api/loader' },
      changeOrigin: true,
    })
  );
  app.use(
    '/api/',
    createProxyMiddleware({
      target: url,
      changeOrigin: true,
    })
  );
};
