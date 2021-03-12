const path = require('path');

module.exports = {
	mode: 'development',
	entry: path.resolve(__dirname, './src/index.tsx'),
	module: {
		rules: [
			{
				test: /\.tsx?$/,
				use: 'ts-loader',
				exclude: /node_modules/,
			},
			{
				test: /\.css$/i,
				use: ["style-loader", "css-loader"],
			},
		],
	},
	resolve: {
		extensions: ['.js', '.ts', '.tsx'],
	},
	output: {
		path: path.resolve(__dirname, './build'),
		filename: 'bundle.js',
	},
};
