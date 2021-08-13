
var express = require('express');
var router = express.Router();

var config = require('../../../config.json');
var dbpool = require('../../lib/dbpool');
var gift = require('../../models/gift');
var area = require('../../models/area');
var session = require('../../models/session');

// ============================================================================

router.get('/', (req, res) => {
    res.render('gift/gen', {
        sess: session.data(req),
        areas: area.idnames(),
    });
});

router.post('/', _A_(async (req, res) => {
    let q = req.body;

    try {
        // check params
        if (!q.expire) throw 'invalid expire date';
        if (q.grpid.length > 4) throw 'grpid [0-9999]';

        let grpid = tonumber(q.grpid);
        let count = tonumber(q.count);
        let aid = tonumber(q.area);
        let reuse = tonumber(q.reuse);
        let expire = Date.fromString(q.expire);
        let memo = q.memo || '';

        if (grpid <= 0 || count <= 0) throw 'error params';

        // rewards
        if (!(q.res_k instanceof Array) || !(q.res_v instanceof Array)) throw 'invalid rewards';

        let rewards = [];
        for (let i = 0; i < q.res_k.length; i++) {
            let res_k = q.res_k[i];
            let res_v = tonumber(q.res_v[i]);

            if (res_k != '' && res_v != 0) {
                rewards.push({ res_k, res_v });
            }
        }
        if (rewards.length == 0) throw 'no rewards';

        // gen codes
        let codes = gift.gen_codes(grpid, count);

        // save
        try {
            let db = dbpool.get('share');

            await db.collection('giftinfo').insertOne({
                _id: grpid,
                area: aid,
                reuse: reuse,
                expire: expire,
                memo: memo,
                rewards: rewards,
                codes: codes,
            });
        } catch (e) {
            if (e.message.match(/^E11000/)) {
                throw 'grpid exists';
            } else {
                throw e;
            }
        }

        // ok
        res.json({ codes }).end();
    } catch (e) {
        res.json({ err: typeof e == 'string' ? e : e.message }).end();
    }
}));

router.post('/load', _A_(async (req, res) => {
    let q = req.body;

    try {
        // check params
        let grpid = tonumber(q.grpid);

        if (grpid <= 0) throw 'error params';

        // load
        let db = dbpool.get('share');

        let doc = await db.collection('giftinfo').findOne({ _id: grpid });
        if (!doc) throw 'grpid NOT found';

        doc.expire = doc.expire.toString();

        res.json(doc).end();
    } catch (e) {
        res.json({ err: typeof e == 'string' ? e : e.message }).end();
    }
}));

// ============================================================================

module.exports = router;
