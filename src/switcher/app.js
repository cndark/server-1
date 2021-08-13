
var express = require('express');
var app = express();

var bodyParser = require('body-parser');
var path = require('path');
var process = require('process');
var morgan = require('morgan');

var config = require('../config.json');

// ====================================
// allow cross site
app.use(function (req, res, next) {
    res.header('Access-Control-Allow-Origin', '*');
    next();
});

// ====================================
// strip /switcher/

app.use((req, res, next) => {
    let m = req.url.match(/^\/switcher(\/.*)/);
    if (m) {
        req.url = m[1];
    }
    next();
});

// ====================================

app.use(bodyParser.urlencoded({ extended: false }));
app.use(morgan('common', {
    skip: (req, res) => req.method == 'HEAD' || res.statusCode < 400 || req.url == '/',
}));

// proxy
app.set('trust proxy', config.common.behind_proxy);

// static
app.use(express.static(path.join(process.env['WORK_DIR'], 'update')));

// ====================================

// routes
app.use('/json', require('./routes/json'));
app.use('/server', require('./routes/server'));
app.use('/notice', require('./routes/notice'));
app.use('/wblist', require('./routes/wblist'));
app.use('/user', require('./routes/user'));

// ====================================

// catch 404
app.use(function (req, res, next) {
    res.status(404).end('Not Found');
});

// error handler
app.use(function (err, req, res, next) {
    res.status(500).end('internal error');
    console.log(err);
});

// ====================================

process.on('uncaughtException', (err) => {
    console.log(err);
});

// ====================================

module.exports = app;
