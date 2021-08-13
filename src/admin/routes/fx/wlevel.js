
var express = require('express');
var router  = express.Router();

var fx      = require('./fx');

// ============================================================================

router.get('/', (req, res) => {
    fx.render(req, res, {
        name:  '关卡停留',

        db:    'stats',
        coll:  'wlevel',

        p0: null,
        filter0: null,

        p1: null,
        filter1: null,

        cols: [
            {type: 'group', key: 'wlv',   text: '关卡进度'},
            {type: 'stat',  key: 'n',     text: '人数', op: {$sum: 1}},
        ],

        hide: ['dr'],
    });
});

// ============================================================================

module.exports = router;
