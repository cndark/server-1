
var express = require('express');
var router = express.Router();

var cluster = require('cluster');
var axios = require('axios');
var md5 = require('md5');
var qs = require('querystring');
var crypto = require('crypto');
var xml2js = require('xml2js');

var config = require('../../config.json');
var dbpool = require('../lib/dbpool');
var gtab = require('../models/gtab');

// ============================================================================

var Auth_SDK_Provider = {

    "soda.internal": async (q, sdkobj) => {
        if (config.common.dev_mode) {
            return;
        } else {
            throw 'not in dev mode';
        }
    },

    "soda.pressure": async (q, sdkobj) => {
        if (config.common.dev_mode) {
            return;
        } else {
            throw 'not in dev mode';
        }
    },

    "soda.ai": async (q, sdkobj) => {
        await dbpool.exec(config.common.db_share, async db => {
            let doc = await db.collection('aibot').findOne({ _id: 1 });
            if (!doc || doc.pwd != q.auth_token) throw 'wrong password';
        });
    },

    "chuxin_and": async (q, sdkobj) => {
        const url = `https://g.chuxinhudong.com/api/index`;

        var form = {
            data: q.auth_token,
            clientIp: q.ip,
        };
        form.sign = md5(`${form.data}${sdkobj.auth_key}`);

        let res2 = await axios.post(url, form);
        let body2 = res2.data;

        if (body2.errno == 1000) {
            return { authid: body2.data.uid };
        } else {
            throw body2;
        };
    },

    "chuxinea_and": async (q, sdkobj) => {
        const url = `https://ea-g.chuxinhudong.com/api/index`;

        var form = {
            data: q.auth_token,
            clientIp: q.ip,
        };
        form.sign = md5(`${form.data}${sdkobj.auth_key}`);

        let res2 = await axios.post(url, form);
        let body2 = res2.data;

        if (body2.errno == 1000) {
            return { authid: body2.data.uid };
        } else {
            throw body2;
        };
    },

    "chuxinea_ios": async (q, sdkobj) => {
        const url = `https://ea-g.chuxinhudong.com/api/index`;

        var form = {
            data: q.auth_token,
            clientIp: q.ip,
        };
        form.sign = md5(`${form.data}${sdkobj.auth_key}`);

        let res2 = await axios.post(url, form);
        let body2 = res2.data;

        if (body2.errno == 1000) {
            return { authid: body2.data.uid };
        } else {
            throw body2;
        };
    },

    "chuxin_h5": async (q, sdkobj) => {
        const url = `https://games.chuxinhudong.com/v1/api/index`;

        let form = q.auth_token;

        // header sign
        let ts = new Date().getTime();
        let algorithm = "CX-HMAC-SHA256";
        let signStr = form + sdkobj.secretId + algorithm + ts;

        var hmac = crypto.createHmac('SHA256', sdkobj.secretKey);
        hmac.update(signStr);
        let sign = hmac.digest('hex');

        let res2 = await axios.post(url, form, {
            headers: {
                "Content-Type": "application/json",
                secretId: sdkobj.secretId,
                algorithm: algorithm,
                signature: sign,
                timestamp: ts,
            }
        });
        let body2 = res2.data;

        if (body2.code == 1000) {
            return {
                authid: body2.data.uid,
                c: {
                    session: body2.data.session,
                    openid: body2.data.open_id,
                },
                s: {
                    openid: body2.data.open_id,
                },
            };
        } else {
            throw body2;
        };
    },

    "wx": async (q, sdkobj) => {
        var url = 'https://api.weixin.qq.com/sns/jscode2session';

        var qs = {
            appid: sdkobj.appid,
            secret: sdkobj.appkey,
            js_code: q.auth_token,
            grant_type: 'authorization_code',
        };

        let res2 = await axios.get(url, {
            params: qs
        });
        let body2 = res2.data;

        if (!body2.errcode || body2.errcode == 0) {
            return { authid: body2.openid };
        } else {
            throw body2;
        }
    },

}

// ============================================================================

router.post('/auth', _A_(async (req, res) => {
    // simple response
    let resp = (b) => {
        res.json({ err: b ? '' : 'failed' }).end();
    }

    // q
    let q = req.body;

    // dyncode auth
    if (q.auth_token.startsWith('soda.bible:') && dyncode.code && q.auth_token.slice(11) == dyncode.code) {
        resp(true);
        return;
    }

    // sdk auth
    let sdkobj = gtab.sdk[q.sdk];
    if (!sdkobj) {
        resp(false);
        return;
    }

    // check provider
    let pvd = q.sdk.match(/^[^-]+/)[0];
    let h = Auth_SDK_Provider[pvd];
    if (!h) {
        resp(false);
        return;
    }

    // ok
    try {
        let ret = await h(q, sdkobj);

        // check ret
        ret = ret || {};
        if (!ret.authid)
            ret.authid = q.auth_id;

        // response
        res.json(ret).end();

        // save auth devid
        save_auth_devid(q.sdk, ret.authid, q.devid).catch(console.error);
    } catch (e) {
        resp(false);
        console.error(`${q.sdk} auth-failed:`, e);
    }
}));

async function save_auth_devid(sdk, authid, devid) {
    try {
        let db = dbpool.get('center');
        await db.collection('acctinfo').updateOne(
            { _id: sdk + '-' + authid },
            { $set: { devid: devid } },
            { upsert: true },
        );
    } catch (e) {
        console.error('saving auth-devid failed:', e.message);
    }
}

// ============================================================================

router.post('/dyncode/update', _A_(async (req, res) => {
    let q = req.body;

    try {
        if (!q.code) throw 'error params';

        cluster.worker.send({ op: 'dyncode.update', body: q.code });
        res.json({}).end();
    } catch (e) {
        res.json({ err: typeof e == 'string' ? e : e.message }).end();
    }
}));

// ============================================================================

var dyncode = {};

cluster.worker.on('message', msg => {
    if (msg.op == 'dyncode.update') {
        dyncode.code = msg.body;
        dyncode.ts = Date.now();
    }
});

// dyncode expiration
setInterval(() => {
    if (!dyncode.code) return;
    if (Date.now() - dyncode.ts < 600 * 1000) return;

    dyncode.code = '';
}, 15 * 1000);

// ============================================================================

module.exports = router;
