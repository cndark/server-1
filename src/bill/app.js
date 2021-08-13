
var express = require('express');
var app = express();

var cluster = require('cluster');
var bodyParser = require('body-parser');
var morgan = require('morgan');

var config = require('../config.json');

// ====================================
// strip /bill/

app.use((req, res, next) => {
    let m = req.url.match(/^\/bill(\/.*)/);
    if (m) {
        req.url = m[1];
    }
    next();
});

// ====================================

app.use(bodyParser.raw({ type: 'application/x-tar' }));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));
app.use(morgan('common', {
    skip: (req, res) => req.method == 'HEAD' || res.statusCode < 400 || req.url == '/',
}));

// proxy
app.set('trust proxy', config.common.behind_proxy);

// ====================================

// routes
app.use('/api', require('./routes/api'));
app.use('/local', require('./routes/local'));

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
