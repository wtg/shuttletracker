const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
    outputDir: '../static',
    baseUrl: '/static/',
    pages: {
        index: {
            entry: 'src/main.ts',
            template: 'public/index.html',
            filename: 'index.html'
        },
        admin: {
            entry: 'src/admin.ts',
            template: 'public/admin.html',
            filename: 'admin.html',
        }
    },
    
    configureWebpack: {
        // plugins: [
        //     new CopyWebpackPlugin([{from: 'src/admin'}], {}),
        // ],
        devServer: {
            proxy: {
                '/vehicles|/routes|/adminMessage|/updates|/stops': {
                    target: 'https://shuttles.rpi.edu',
                    changeOrigin: true,
                    cookieDomainRewrite: '',
                },
            }
        }
    },
    chainWebpack: config => {
        config.optimization.delete('splitChunks');
    }
};