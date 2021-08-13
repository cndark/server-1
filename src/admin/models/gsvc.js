
var axios = require('axios');

var dbpool = require('../lib/dbpool');
var utils  = require('../lib/utils');
var area   = require('./area');

// ============================================================================

async function run(operator, q, areaid, svrstr) {
    // get area conf
    let a = area.find(areaid);
    if (!a) throw 'area NOT found';

    // parse svrstr
    let ids = utils.parse_svrstr(a.conf, svrstr);
    if (ids.length == 0) throw 'no valid svr found';

    // run svc
    let r = {};
    for (let id of ids) {
        let svr = `game${id}`;

        // get gs object
        let gs = a.conf.games[svr];
        if (!gs) continue;

        // run
        try {
            let {data} = await axios.post(`http://${gs.svc}/gm/service`, q);
            r = data;
        } catch(e) {
            console.error('run gsvc failed:', typeof e == 'string' ? e : e.message);
        }
    }

    // write gm log
    try {
        if (!r.err || r.err == '')
            await write_gm_log(operator, q, areaid, svrstr);
    } catch(e) {
        console.error('write GM-LOG failed:', typeof e == 'string' ? e : e.message);
    }

    // ok
    return r;
}

async function write_gm_log(operator, q, areaid, svrstr) {
    if (!q.key || !q.key.startsWith("w.")) return;

    let db = dbpool.get('share');

    await db.collection("gmlog").insertOne({
        operator: operator,
        area:     areaid,
        tarsvr:   svrstr,
        userid:   q.plrid,
        op:       q,
        ts:       new Date(),
    });
}

// ============================================================================

module.exports = {
    run: run,
}
