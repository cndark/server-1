
var express = require('express');
var router  = express.Router();

var gtab    = require('../../models/gtab');

// ============================================================================

router.use('/gen',      require('./gen'));
router.use('/usage',    require('./usage'));

// ============================================================================

router.post('/gtab', (req, res) => {
    var t = gtab[req.body.key];
    if (!t) t = [["*", "*"]];

    res.json(t).end();
});

// ============================================================================

module.exports = router;
