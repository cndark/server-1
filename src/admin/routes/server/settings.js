
var express = require('express');
var router  = express.Router();

var axios = require('axios');

var config  = require('../../../config.json');
var dbpool  = require('../../lib/dbpool');
var area    = require('../../models/area');
var session = require('../../models/session');

// ============================================================================

router.get('/', _A_(async (req, res) => {
    res.render('svr/settings', {
        sess:  session.data(req),
        areas: area.idnames(),
    });
}));

router.post('/load', _A_(async (req, res) => {
    let q = req.body;

    let aid = tonumber(q.area);

    try {
        // area
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // load
        let db = dbpool.get('share');

        let doc = await db.collection('settings').findOne({_id: aid});
        if (!doc) throw 'not found';

        res.json(doc).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

router.post('/save', _A_(async (req, res) => {
    let q = req.body;

    let aid             = tonumber(q.area);
    let closereg_hwater = tonumber(q.closereg_hwater);
    let closereg_limit  = tonumber(q.closereg_limit);
    let opennew_mode    = q.opennew_mode;
    let opennew_itv     = tonumber(q.opennew_itv);

    try {
        // area
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // set-doc
        let doc = {};

        // closereg hwater
        if (closereg_hwater >= 0) {
            doc.closereg_hwater = closereg_hwater;
        }

        // closereg limit
        if (closereg_limit >= 0) {
            doc.closereg_limit = closereg_limit;
        }

        // opennew mode
        if (opennew_mode) {
            doc.opennew_mode = opennew_mode;
        }

        // opennew itv
        if (opennew_itv >= 0) {
            doc.opennew_itv = opennew_itv;
        }

        // check doc
        if (Object.keys(doc).length == 0) throw 'no update specified';

        // update
        let db = dbpool.get('share');
        await db.collection('settings').updateOne({_id: aid}, {$set: doc}, {upsert: true});

        // notify agent
        await axios.post(
            `http://${config.agent.ip}:${config.agent.port}/settings/update`,
        );

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
