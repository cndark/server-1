
var express = require('express');
var router  = express.Router();

var settings = require('../models/settings');

// ============================================================================

router.post('/update', _A_(async (req, res) => {
    try {
        await settings.load();

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
