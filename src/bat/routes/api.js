
var express = require('express');
var router = express.Router();
var battle = require('../calcbattle/index');

// ============================================================================

router.post('/calc', function (req, res) {
    var q = req.body;

    var r = battle.startBattle(q);

    res.json(r).end();
});

// ============================================================================

module.exports = router;
