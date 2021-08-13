
var express = require('express');
var router  = express.Router();

var session = require('../models/session');

// ============================================================================

router.get('/', (req, res) => {
    // already logged in
    if (req.session.user) {
        res.redirect('/main');
        return;
    }

    // show login page
    res.render('login', {login_err: ""});
});

router.post('/', _A_(async (req, res) => {
    let q = req.body;

    try {
        if (! await session.auth(req, q.user, q.pwd)) throw 'failed';

        res.redirect('/main');
    } catch {
        res.render('login', {login_err: "Invalid account or password"});
    }
}));

// ============================================================================

module.exports = router;
