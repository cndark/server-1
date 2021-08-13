
var express = require('express');
var router = express.Router();

var gtab = require('../../models/gtab');

// ============================================================================

router.use('/bill_items', require('./bill_items'));
router.use('/bill_first', require('./bill_first'));
router.use('/vip', require('./vip'));
router.use('/svr_roll', require('./svr_roll'));
router.use('/wlevel', require('./wlevel'));
router.use('/tutorial', require('./tutorial'));

// ============================================================================

router.post('/gtab', (req, res) => {
    var t = gtab[req.body.key];
    if (!t) t = [["*", "*"]];

    res.json(t).end();
});

// ============================================================================

module.exports = router;
