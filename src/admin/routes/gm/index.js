
var express = require('express');
var router  = express.Router();

var gtab    = require('../../models/gtab');

// ============================================================================

router.use('/tool',     require('./tool'));
router.use('/user',     require('./user'));
router.use('/guild',    require('./guild'));
router.use('/tmail',    require('./tmail'));
router.use('/log',      require('./log'));

// ============================================================================

router.post('/gtab', (req, res) => {
    var t = gtab[req.body.key];
    if (!t) t = [["*", "*"]];

    res.json(t).end();
});

// ============================================================================

module.exports = router;
