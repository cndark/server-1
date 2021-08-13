
var express = require('express');
var router  = express.Router();

var gtab    = require('../../models/gtab');

// ============================================================================

router.use('/user',   require('./user'));

// ============================================================================

router.post('/gtab', function (req, res) {
    var t = gtab[req.body.key];
    if (!t) t = [["*", "*"]];

    res.json(t).end();
});

// ============================================================================

module.exports = router;
