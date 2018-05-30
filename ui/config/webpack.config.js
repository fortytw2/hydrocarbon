const autoprefixer = require("autoprefixer");
const path = require("path");
const webpack = require("webpack");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CaseSensitivePathsPlugin = require("case-sensitive-paths-webpack-plugin");
const CleanWebpackPlugin = require("clean-webpack-plugin");

const getClientEnvironment = require("./env");
const paths = require("./paths");

const jsxCompilationOptions = {
  compilationOptions: {},
  types: "typescript",
  aliases: {},
  optimization: {
    rewriteIdents: true,
    mergeDeclarations: true,
    removeUnusedStyles: true,
    conflictResolution: true,
    enabled: false
  }
};

// Webpack uses `publicPath` to determine where the app is being served from.
// In development, we always serve from the root. This makes config easier.
const publicPath = "/";
// `publicUrl` is just like `publicPath`, but we will provide it to our app
// as %PUBLIC_URL% in `index.html` and `process.env.PUBLIC_URL` in JavaScript.
// Omit trailing slash as %PUBLIC_PATH%/xyz looks better than %PUBLIC_PATH%xyz.
const publicUrl = "/";
// Get environment variables to inject into our app.
const env = getClientEnvironment(publicUrl);

// This is the development configuration.
// It is focused on developer experience and fast rebuilds.
// The production configuration is different and lives in a separate file.
module.exports = {
  // You may want 'eval' instead if you prefer to see the compiled output in DevTools.
  // See the discussion in https://github.com/facebookincubator/create-react-app/issues/343.
  devtool: "cheap-module-source-map",
  // These are the "entry points" to our application.
  // This means they will be the "root" imports that are included in JS bundle.
  // The first two entry points enable "hot" CSS and auto-refreshes for JS.
  entry: [
    // We ship a few polyfills by default:
    require.resolve("./polyfills"),
    // Include an alternative client for WebpackDevServer. A client's job is to
    // connect to WebpackDevServer by a socket and get notified about changes.
    // When you save a file, the client will either apply hot updates (in case
    // of CSS changes), or refresh the page (in case of JS changes). When you
    // make a syntax error, this client will display a syntax error overlay.
    // Note: instead of the default WebpackDevServer client, we use a custom one
    // to bring better experience for Create React App users. You can replace
    // the line below with these two lines if you prefer the stock client:
    // require.resolve("webpack-dev-server/client") + "?/",
    // require.resolve("webpack/hot/dev-server"),
    // Finally, this is your app's code:
    paths.appIndexJs
  ],
  output: {
    pathinfo: true,
    // This does not produce a real file. It's just the virtual path that is
    // served by WebpackDevServer in development. This is the JS bundle
    // containing code from all our entry points, and the Webpack runtime.
    filename: "static/hydrocarbon.[hash:8].js",
    // This is the URL that app is served from. We use "/" in development.
    publicPath: publicPath,
    // Point sourcemap entries to original disk location (format as URL on Windows)
    devtoolModuleFilenameTemplate: info =>
      path.resolve(info.absoluteResourcePath).replace(/\\/g, "/")
  },
  resolve: {
    // This allows you to set a fallback for where Webpack should look for modules.
    // We placed these paths second because we want `node_modules` to "win"
    // if there are any conflicts. This matches Node resolution mechanism.
    // https://github.com/facebookincubator/create-react-app/issues/253
    modules: ["node_modules", paths.appNodeModules].concat(
      // It is guaranteed to exist because we tweak it in `env.js`
      process.env.NODE_PATH.split(path.delimiter).filter(Boolean)
    ),
    alias: {
      "@": paths.appSrc
    },
    extensions: [".js", ".json", ".jsx", ".ts", ".tsx"]
  },
  module: {
    strictExportPresence: true,
    rules: [
      // TODO: Disable require.ensure as it's not a standard language feature.
      // We are waiting for https://github.com/facebookincubator/create-react-app/issues/2176.
      // { parser: { requireEnsure: false } },

      {
        // "oneOf" will traverse all following loaders until one will
        // match the requirements. When no loader matches it will fall
        // back to the "file" loader at the end of the loader list.
        oneOf: [
          // "url" loader works like "file" loader except that it embeds assets
          // smaller than specified limit in bytes as data URLs to avoid requests.
          // A missing `test` is equivalent to a match.
          {
            test: [/\.bmp$/, /\.gif$/, /\.jpe?g$/, /\.png$/],
            loader: require.resolve("url-loader"),
            options: {
              limit: 10000,
              name: "static/media/[name].[hash:8].[ext]"
            }
          },

          // // "postcss" loader applies autoprefixer to our CSS.
          // // "css" loader resolves paths in CSS and adds assets as dependencies.
          // // "style" loader turns CSS into JS modules that inject <style> tags.
          // // In production, we use a plugin to extract that CSS to a file, but
          // // in development "style" loader enables hot editing of CSS.
          {
            test: /\.css$/,
            exclude: /node_modules/,
            use: [
              MiniCssExtractPlugin.loader,
              {
                loader: "css-loader",
                options: {
                  modules: true,
                  importLoaders: 1,
                  localIdentName: "hc__[hash:base64:8]"
                }
              },
              {
                loader: "postcss-loader",
                options: {
                  sourceMap: true,
                  ident: "postcss",
                  plugins: loader => [
                    require("autoprefixer")(),
                    require("cssnano")()
                  ]
                }
              }
            ]
          },

          {
            test: /\.js(x?)$/,
            exclude: /node_modules/,
            use: [
              {
                loader: require.resolve("babel-loader"),
                options: {
                  presets: [
                    "@babel/preset-env",
                    ["@babel/preset-stage-1", { decoratorsLegacy: true }]
                  ],
                  plugins: [
                    ["@babel/plugin-proposal-decorators", { legacy: true }],
                    [
                      "@babel/plugin-proposal-class-properties",
                      { loose: true }
                    ],
                    ["babel-plugin-transform-react-jsx", { pragma: "h" }],
                    [
                      "babel-plugin-jsx-pragmatic",
                      {
                        module: "preact",
                        export: "h",
                        import: "h"
                      }
                    ]
                  ],
                  cacheDirectory: true,
                  compact: true
                }
              }
            ]
          },

          // "file" loader makes sure those assets get served by WebpackDevServer.
          // When you `import` an asset, you get its (virtual) filename.
          // In production, they would get copied to the `build` folder.
          // This loader doesn't use a "test" so it will catch all modules
          // that fall through the other loaders.
          {
            // Exclude `js` files to keep "css" loader working as it injects
            // it's runtime that would otherwise processed through "file" loader.
            // Also exclude `html` and `json` extensions so they get processed
            // by webpacks internal loaders.
            exclude: [/\.[j|t]s(x?)$/, /\.html$/, /\.json$/],
            loader: require.resolve("file-loader"),
            options: {
              name: "media/[name].[hash:8].[ext]"
            }
          }
        ]
      }
      // ** STOP ** Are you adding a new loader?
      // Make sure to add the new loader(s) before the "file" loader.
    ]
  },
  plugins: [
    new CleanWebpackPlugin(["dist"], {
      root: paths.root
    }),
    new MiniCssExtractPlugin({
      // Options similar to the same options in webpackOptions.output
      // both options are optional
      filename: "static/hydrocarbon.[contenthash:8].css"
    }),
    // Generates an `index.html` file with the <script> injected.
    new HtmlWebpackPlugin({
      inject: true,
      template: paths.appHtml,
      filename: "index.html"
    }),
    // Add module names to factory functions so they appear in browser profiler.
    new webpack.NamedModulesPlugin(),
    // Makes some environment variables available to the JS code, for example:
    // if (process.env.NODE_ENV === 'development') { ... }. See `./env.js`.
    new webpack.DefinePlugin(env.stringified),
    // Watcher doesn't work well if you mistype casing in a path so we use
    // a plugin that prints an error when you attempt to do this.
    // See https://github.com/facebookincubator/create-react-app/issues/240
    new CaseSensitivePathsPlugin(),
    // Moment.js is an extremely popular library that bundles large locale files
    // by default due to how Webpack interprets its code. This is a practical
    // solution that requires the user to opt into importing specific locales.
    // https://github.com/jmblog/how-to-optimize-momentjs-with-webpack
    // You can remove this if you don't use Moment.js:
    new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/)
  ],
  // Some libraries import Node modules but don't use them in the browser.
  // Tell Webpack to provide empty mocks for them so importing them works.
  node: {
    dgram: "empty",
    fs: "empty",
    net: "empty",
    tls: "empty",
    child_process: "empty"
  },
  // Turn off performance hints during development because we don't do any
  // splitting or minification in interest of speed. These warnings become
  // cumbersome.
  performance: {
    hints: false
  }
};
