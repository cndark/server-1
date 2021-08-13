
var express = require('express');
var router  = express.Router();

var gtab    = require('../../models/gtab');

// ============================================================================

router.use('/status',   require('./status'));
router.use('/onoff',    require('./onoff'));
router.use('/opennew',  require('./opennew'));
router.use('/notice',   require('./notice'));
router.use('/wblist',   require('./wblist'));
router.use('/gdata',    require('./gdata'));
router.use('/settings', require('./settings'));

// ============================================================================

router.post('/gtab', function (req, res) {
    var t = gtab[req.body.key];
    if (!t) t = [["*", "*"]];

    res.json(t).end();
});

// ============================================================================

module.exports = router;
