
// ============================================================================

function run(app, port, host) {
    let server = require('http').createServer(app);

    server.on('error', (err) => {
        switch (err.code) {
            case 'EACCES':
                console.error(`port ${port} requires elevated privileges`);
                process.exit(1);
                break;

            case 'EADDRINUSE':
                console.error(`port ${port} is already in use`);
                process.exit(1);
                break;

            default:
                throw err;
        }
    });

    server.on('listening', () => {
        console.info(`Listening on ${host ? host : ''}:${port}`);
    });

    server.listen(port, host);
}

// ============================================================================

module.exports = {
    run: run,
};
