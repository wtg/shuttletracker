module.exports = {
    outputDir: '../static',
    
    configureWebpack: {
        
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