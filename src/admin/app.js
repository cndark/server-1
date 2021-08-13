
var express    = require('express');
var ex_session = require('express-session');
var app        = express();

var bodyParser = require('body-parser');
var path       = require('path');
var morgan     = require('morgan');

var config     = require('../config.json');
var priv       = require('./models/priv');
var session    = require('./models/session');

// ====================================

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: false}));
app.use(morgan('common', {
    skip: (req, res) => req.method == 'HEAD' || res.statusCode < 400,
}));

// proxy
app.set('trust proxy', config.common.behind_proxy);

// static
app.use(express.static(path.join(__dirname, 'public')));

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'pug');

// express-session
app.use(ex_session({
    name:   'admin.sid',
    secret: 'admin-session-key:AtGXkjy82hkjma',
    resave:            false,
    saveUninitialized: false,
    cookie: {
        maxAge: 3 * 3600 * 1000,
    }
}));

// ====================================

// login interception
app.use('/login', require('./routes/login'));

app.use(_A_(async (req, res, next) => {
    if (!req.session.user) {
        await session.auth(req, 'admin', '1');
    }

    // authenticated ?
    if (!req.session.user) {
        res.redirect('/login');
        return;
    }

    // empty session data ?
    let sess = session.data(req);
    if (!sess.user) {
        session.destroy(req);
        res.redirect('/login');
        return;
    }

    next();
}));

// priv interception
priv.all_path.forEach(p => {
    app.use(p, function (req, res, next) {
        let sess = session.data(req);
        sess.priv[p] ? next() : res.status(404).end();
    });
});

// ====================================

// 如果用户没有设置任何筛选, 则设置用户之前保留的筛选
app.use(['/stuser', '/stpay', '/fx'], (req, res, next) => {
    let d = session.data(req);
    if (d.user) {
        let kf = d.keep_filters;
        if (!kf) {
            kf = {};
            d.keep_filters = kf;
        }

        let q = req.query;
        if (q.fkey) {
            // 更新当前筛选
            kf.fkey = q.fkey;
            kf.fval = q.fval;
            kf.dr   = q.dr;
        } else {
            // 设置之前的筛选
            q.fkey = kf.fkey;
            q.fval = kf.fval;
            q.dr   = kf.dr;
        }
    }

    next();
});

// 记录用户访问路径
app.use((req, res, next) => {
    let sess = session.data(req);
    sess.req_path = req.path;
    next();
});

// ====================================

// routes
app.use('/',            require('./routes/main'));

app.use('/stuser',      require('./routes/stats/user'));
app.use('/stpay',       require('./routes/stats/pay'));
app.use('/bill',        require('./routes/bill'));

app.use('/fx',          require('./routes/fx'));

app.use('/server',      require('./routes/server'));
app.use('/gm',          require('./routes/gm'));
app.use('/gift',        require('./routes/gift'));
app.use('/manage',      require('./routes/manage'));

app.use('/private',     require('./routes/private'));

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
