
var express = require('express');
var router  = express.Router();

var axios   = require('axios');

var dbpool  = require('../../lib/dbpool');
var area    = require('../../models/area');
var gtab    = require('../../models/gtab');
var gsvc    = require('../../models/gsvc');
var session = require('../../models/session');

// ============================================================================

router.get('/', function(req, res) {
    res.render('svr/notice', {
        sess:      session.data(req),
        areas:     area.idnames(),
        languages: gtab.lang,
    });
});

router.post('/load', _A_(async (req, res) => {
    var q = req.body;

    q.area = tonumber(q.area)

    try {
        // get area conf
        let a = area.find(q.area);
        if (!a) throw 'area NOT found';

        // check lang
        if (!q.lang) throw 'invalid lang';

        // load
        let doc = await dbpool.exec(a.conf.common.db_center, async db => {
            let doc = await db.collection('notice').findOne({_id: q.lang});
            if (doc) {
                doc.d_start = doc.d_start.toString();
                doc.d_end   = doc.d_end.toString();
            }
            return doc;
        });

        doc ? res.json(doc).end() : res.json({err: "notice not found"}).end();
    } catch(e) {
        res.json({err: 'notice load failed'}).end();
        console.error('notice load failed:', e);
    }
}));

router.post('/save', _A_(async (req, res) => {
    var q = req.body;

    if (!q.area || !q.lang || !q.d_start || !q.d_end || !q.title || !q.content) {
        res.json({err: "error params"}).end();
        return;
    }

    q.area    = tonumber(q.area);
    q.d_start = Date.fromString(q.d_start);
    q.d_end   = Date.fromString(q.d_end);

    try {
        // get area conf
        let a = area.find(q.area);
        if (!a) throw 'area NOT found';

        // save
        await dbpool.exec(a.conf.common.db_center, async db => {
            await db.collection('notice').replaceOne(
                {_id: q.lang},
                {
                    d_start: q.d_start,
                    d_end:   q.d_end,
                    title:   q.title,
                    content: q.content,
                },
                {upsert: true},
            );
        });

        // notify switcher
        await axios.get(
            `http://${a.conf.switcher.ip}:${a.conf.switcher.port}/notice/update?token=${a.conf.switcher.token}`,
        );

        res.json({}).end();
    } catch(e) {
        res.json({err: 'notice save failed'}).end();
        console.error('notice save failed:', e);
    }
}));

// ============================================================================

router.post('/lamp', _A_(async (req, res) => {
    let sess = session.data(req);
    let q    = req.body;

    if (!q.area || !q.target || !q.content) {
        res.json({err: "error params"}).end();
        return;
    }

    q.key  = 'w.lamp';
    q.area = tonumber(q.area);

    try {
        let r = await gsvc.run(sess.user, q, q.area, q.target);
        res.json(r).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
