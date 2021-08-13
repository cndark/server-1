var express = require('express');
var router = express.Router();

var glog = require('../../models/glog');
var token = require('../../models/token');

function check_param(q) {
    // check param
    if (!q.op || !q.area || !q.svr) {
        return false
    }

    now = new Date();
    q.ts = now.toString();

    return true
}

// ============================================================================

router.post('/', _A_(async (req, res) => {
    let q = req.body;
    if (!token.check_token(q)) {
        res.end("token err");
        return
    }
    delete q.token;

    if (!check_param(q)) {
        res.end("param err");
        return
    }

    try {
        await glog.send(JSON.stringify(q));

        res.end("0");
    } catch (e) {
        res.end("1");
    }
}));

// ============================================================================

module.exports = router;
