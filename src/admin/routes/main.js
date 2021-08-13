
var express = require('express');
var router  = express.Router();
var session = require('../models/session');

// ============================================================================

router.get('/', (req, res) => {
    res.redirect('/main');
});

router.get('/main', (req, res) => {
    res.render('main', {
        sess: session.data(req),
    });
});

router.get('/logout', (req, res) => {
    session.destroy(req);
    res.redirect('/login');
});

// ============================================================================

module.exports = router;
