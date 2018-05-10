var ExtractTextPlugin = require('extract-text-webpack-plugin')
var CopyWebpackPlugin = require('copy-webpack-plugin')
var Webpack = require('webpack')

module.exports = {
  entry: {
    app: [
      './ui/index.js'
    ]
  },

  output: {
    path: require('path').resolve(__dirname, 'assets'),
    filename: 'proxy.js'
  },

  devtool: 'source-map',

  resolve: {
    extensions: [
      '',
      '.js'
    ]
  },

  module: {
    loaders: [
      {
        test: /\.js?$/,
        exclude: /node_modules/,
        loader: 'babel',
        query: {
          presets: [
            'es2015'
          ]
        }
      },
      {
        test: /\.less$/,
        loader: ExtractTextPlugin.extract(
          'css?sourceMap!less?sourceMap'
        )
      },
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file?name=[name].[ext]&mimetype=application/font-woff'
      },
      {
        test: /\.otf(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file?name=[name].[ext]&mimetype=application/x-font-opentype'
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file?name=[name].[ext]&mimetype=image/svg+xml'
      },
      {
        test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file?name=[name].[ext]&mimetype=application/octet-stream'
      },
      {
        test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file?name=[name].[ext]&mimetype=application/vnd.ms-fontobject'
      }
    ]
  },

  vue: {
    loaders: {
      css: ExtractTextPlugin.extract('css!less')
    }
  },

  plugins: [
    new ExtractTextPlugin(
      'proxy.css'
    ),
    new CopyWebpackPlugin([{
      from: 'ui/images',
      to: ''
    }]),
    new Webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(
        process.env.NODE_ENV || 'development'
      )
    }),
    new Webpack.optimize.OccurenceOrderPlugin(),
    new Webpack.optimize.DedupePlugin()
  ]
}

if (process.env.NODE_ENV === 'production') {
  module.exports.devtool = '#source-map'

  module.exports.plugins = (module.exports.plugins || []).concat([
    new Webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: false
      },
      output: {
        semicolons: false
      }
    })
  ])
}
