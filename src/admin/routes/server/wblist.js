
var express = require('express');
var router  = express.Router();

var axios = require('axios');

var dbpool  = require('../../lib/dbpool');
var area    = require('../../models/area');
var session = require('../../models/session');

// ============================================================================

router.get('/', (req, res) => {
    res.render('svr/wblist', {
        sess:  session.data(req),
        areas: area.idnames(),
    });
});

router.get('/load', _A_(async (req, res) => {
    let q = req.query;

    q.area = tonumber(q.area);

    try {
        // get area conf
        let a = area.find(q.area);
        if (!a) throw 'area NOT found';

        // load
        let doc = await dbpool.exec(a.conf.common.db_center, async db => {
            let doc = await db.collection('wblist').findOne({_id: 1});
            if (!doc) {
                doc = {w_ips:[], w_devices:[], b_ips:[], b_devices:[]};
            }
            return doc;
        });

        res.json(doc).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

router.post('/save', _A_(async (req, res) => {
    let q = req.body;

    q.area = tonumber(q.area);

    try {
        // get area conf
        let a = area.find(q.area);
        if (!a) throw 'area NOT found';

        // save
        let doc = {
            w_ips:     req.body.w_ips.split(/[\r\n]+/),
            w_devices: req.body.w_devices.split(/[\r\n]+/),
            b_ips:     req.body.b_ips.split(/[\r\n]+/),
            b_devices: req.body.b_devices.split(/[\r\n]+/),
        };

        await dbpool.exec(a.conf.common.db_center, async db => {
            await db.collection('wblist').updateOne({_id: 1}, {$set: doc}, {upsert: true});
        });

        // notify switcher
        await axios.get(
            `http://${a.conf.switcher.ip}:${a.conf.switcher.port}/wblist/update?token=${a.conf.switcher.token}`,
        );

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
