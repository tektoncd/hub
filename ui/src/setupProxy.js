const {createProxyMiddleware} = require('http-proxy-middleware');

module.exports = function (app) {
  app.use(
    createProxyMiddleware(['/auth', '/user'], {
      target: 'http://localhost:4200',
      changeOrigin: true
    })
  );
  app.use(
    createProxyMiddleware(['/resource/**/rating'], {
      target: 'http://localhost:8000',
      changeOrigin: true
    })
  );
};
