export default function(config, env, helpers) {
  if (env.production) {
    config.output.publicPath = "/static/";
  }

  if (!env.production) {
    config.devServer.proxy = [
      {
        path: "/v1/**",
        changeOrigin: true,
        changeHost: true,
        secure: false,
        target: "http://localhost:8081"
      }
    ];
  }
}
