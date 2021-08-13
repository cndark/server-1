
var express    = require('express');
var app        = express();

var bodyParser = require('body-parser');
var morgan     = require('morgan');

// ====================================

app.use(bodyParser.json());
app.use(morgan('common', {
    skip: (req, res) => req.method == 'HEAD' || res.statusCode < 400,
}));

// ====================================

// routes
app.use('/api',   require('./routes/api'));

// ====================================

// catch 404
app.use(function(req, res, next) {
    res.status(404).end('Not Found');
});

// error handler
app.use(function(err, req, res, next) {
    res.status(500).end('internal error');
    console.log(err);
});

// ====================================

process.on('uncaughtException', (err) => {
    console.log(err);
});

// ====================================

module.exports = app;
