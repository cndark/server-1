
var express = require('express');
var router = express.Router();

var config = require('../../config.json');
var dbpool = require('../lib/dbpool');
var token = require('../models/token');

// ============================================================================

router.post('/get', _A_(async (req, res) => {
    let q = req.body;

    // check
    if (!q.auth_id || !q.sdk || !q.devid) {
        res.json({ err: 'err-args' }).end();
        return;
    }

    // gen token
    res.json({ token: token.encode(q.auth_id, q.sdk, q.devid) }).end();
}));

router.post('/auth', _A_(async (req, res) => {
    let q = req.body;

    // verify token
    do {
        if (!q.token) break;

        let obj = token.decode(q.token);
        if (!obj) break;
        if (q.auth_id == "-" || obj.auth_id != q.auth_id) break;
        if (obj.sdk != q.sdk || obj.devid != q.devid) break;
        if (obj.expire * 1000 - Date.now() < 0) break;

        // check devid against last sdk-auth-devid
        if (! await check_devid(q.sdk, q.auth_id, q.devid)) break;

        // ok
        res.json({}).end();
        return;

    } while (true);

    // failed
    res.json({ err: 'err-token' }).end();
}));

async function check_devid(sdk, authid, devid) {
    try {
        let db = dbpool.get('center');
        let doc = await db.collection('acctinfo').findOne({ _id: sdk + '-' + authid });
        if (!doc || devid != doc.devid) throw 'failed';

        return true;
    } catch (e) {
        return false;
    }
}

// ============================================================================

module.exports = router;
