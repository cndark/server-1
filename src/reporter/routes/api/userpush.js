
var express = require('express');
var router = express.Router();

var cluster = require('cluster');
var crypto = require('crypto');
var axios = require('axios');

var config = require('../../../config.json');

// ============================================================================

router.post('/wx_h5', _A_(async (req, res) => {
    // token is needed
    if (req.query.token != config.reporter.token) {
        res.status(404).end();
        return;
    }

    // q
    let q = req.body;

    try {
        // check params
        if (!q.openid || !q.tpl) throw 'error params';

        if (!q.vals) {
            q.vals = [];
        } else if (!(q.vals instanceof Array)) {
            q.vals = [q.vals];
        }

        // get access token
        let atoken = await get_cx_atoken();

        // send msg
        await wx_push(atoken, q);

        // ok
        res.end('0');
    } catch (e) {
        res.end('1');
        console.log('userpush failed:', typeof e == 'string' ? e : e.message);
    }
}));

async function get_cx_atoken() {
    const url = `https://games.chuxinhudong.com/v1/api/accessToken`;
    const secret_id = 'f4HQi9jJ2V2t';
    const secret_Key = 'dGRnKu7bsH0Y2zzEuy5G';
    const algo = 'CX-HMAC-SHA256';

    let body = {
        game: 'mmzz',
        platform: 'wechat',
    };
    let ts = new Date().unix();

    // sign
    var hmac = crypto.createHmac('SHA256', secret_Key);
    hmac.update(`${secret_id}${algo}${ts}`);
    let sign = hmac.digest('hex');

    // request
    let { data } = await axios.get(url, {
        params: body,
        headers: {
            secretId: secret_id,
            algorithm: algo,
            signature: sign,
            timestamp: ts,
        }
    });

    // result
    if (data.code != 1000) throw `${data.code} : ${data.msg}`;

    return data.data.access_token;
}

async function wx_push(atoken, q) {
    const url = 'https://api.weixin.qq.com/cgi-bin/message/subscribe/send';

    // find template
    let templates = require('../../models/userpush');
    let tpl = templates[q.tpl];
    if (!tpl) throw 'error tpl';

    // clone data
    let d = {};
    for (let k in tpl.data) {
        d[k] = { value: tpl.data[k] };
    }

    // replace params
    tpl.keys.forEach((k, i) => {
        d[k].value = q.vals[i] || '';
    });

    // body
    let body = {
        touser: q.openid,
        template_id: tpl.id,
        data: d,
        miniprogram_state: 'formal',
    };

    // push
    let { data } = await axios.post(url, body, {
        headers: {
            "Content-Type": "application/json",
        },
        params: {
            access_token: atoken,
        }
    });

    // check
    if (data.errcode != 0) throw `${data.errcode} : ${data.errmsg}`;
}

// ============================================================================

cluster.worker.on('message', msg => {
    if (msg.op == 'reload') {
        delete require.cache[require.resolve('../../models/userpush')];
    }
});

// ============================================================================

module.exports = router;
