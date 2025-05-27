module.exports = {
  transpileDependencies: [],
  devServer: {
    port: 3000,
    open: true,
    proxy: {
      '^/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
}