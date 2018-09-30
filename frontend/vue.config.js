const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
    outputDir: '../static',
    
    configureWebpack: {
        plugins: [
            new CopyWebpackPlugin([{from: 'src/admin'}], {}),
        ],
        devServer: {
            proxy: {
                '/vehicles|/routes|/adminMessage|/updates|/stops': {
                    target: 'https://shuttles.rpi.edu',
                    changeOrigin: true,
                    cookieDomainRewrite: '',
                },
            }
        }
    }
};