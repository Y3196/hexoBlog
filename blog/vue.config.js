module.exports = {
  publicPath: '/',
  transpileDependencies: ["vuetify"],
  devServer: {
      historyApiFallback: true,
      open: true,
    proxy: {
      "/api": {
       target: "http://localhost:8080/",
      // target: "http://192.168.0.108:8080/",
        changeOrigin: true,
        pathRewrite: {
          "^/api": ""
        },
          onError(err) {
              console.error('Proxy error:', err);
          },
        timeout: 60000, // 增加超时时间为 60 秒
          logLevel: "debug", // 添加此行以启用调试日志
      }
    },
    disableHostCheck: true
  },
  productionSourceMap: false,
  css: {
    extract: true,
    sourceMap: false
  }
};
