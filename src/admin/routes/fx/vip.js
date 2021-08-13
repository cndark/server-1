
var express = require('express');
var router  = express.Router();

var fx      = require('./fx');

// ============================================================================

router.get('/', (req, res) => {
    fx.render(req, res, {
        name:  'Vip 统计',

        db:    'share',
        coll:  'userinfo',

        p0: null,
        filter0: null,

        p1: null,
        filter1: null,

        cols: [
            {type: 'group', key: 'svr',  text: '区服', selectable: true},
            {type: 'group', key: 'vip',  text: 'Vip 等级', selectable: true},
            {type: 'stat',  key: 'n',    text: 'Vip 人数', op: {$sum: 1}},
        ],

        hide: ['fkey', 'dr'],
    });
});

// ============================================================================

module.exports = router;
