const createProxyMiddleware = require('http-proxy-middleware');
const dev = !process.env.NODE_ENV || process.env.NODE_ENV === 'development';
const dockerhost = dev
  ? process.env.REACT_APP_DOCKER_PATH || '192.168.99.100'
  : 'localhost';
const port = 8080;
const url = `http://${dockerhost}:${port}`;

module.exports = function (app) {
  app.use(
    '/api/auth/',
    createProxyMiddleware({
      target: url,
      changeOrigin: true,
    })
  );
  app.use(
    '/api/profile/',
    createProxyMiddleware({
      target: url,
      changeOrigin: true,
    })
  );
  app.use(
    '/api/passwd/',
    createProxyMiddleware({
      target: url,
      changeOrigin: true,
    })
  );
  app.use(
    '/api/search/',
    createProxyMiddleware({
      target: url,
      // pathRewrite: { '^/api/search': '' },
      changeOrigin: true,
    })
  );
  app.use(
    '/api/movies/',
    createProxyMiddleware({
      target: url,
      // pathRewrite: { '^/api/movies': '' },
      changeOrigin: true,
    })
  );
  app.use(
    '/api/storage/load',
    createProxyMiddleware({
      target: url,
      changeOrigin: true,
    })
  );
  app.use(
    '/api/loader/',
    createProxyMiddleware({
      target: url,
      changeOrigin: true,
    })
  );
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
      // pathRewrite: { '^/api': '' },
      changeOrigin: true,
    })
  );
};
