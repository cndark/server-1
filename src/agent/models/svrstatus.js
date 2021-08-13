
var axios = require('axios');

var dbpool = require('../lib/dbpool');
var area   = require('./area');

// ============================================================================

var status_onoff = {
    // aid: {
    //     gameid: {
    //         st: 'off|on|starting|stopping',
    //         ts: unix,
    //     }
    // },
};

async function fetch_status_onoff() {
    for (let aid of area.ids()) {
        let a = area.find(aid);

        let old_lst = status_onoff[aid] || {};

        try {
            let {data} = await axios.get(`http://${a.conf.switcher.ip}:${a.conf.switcher.port}/server/list/games`);

            let new_lst = {};

            for (let k in data) {
                let e = data[k];

                let new_st = e.on ? 'on' : 'off';
                let new_ts = Date.now() / 1000;

                let old_obj = old_lst[e.id];
                if (old_obj) {
                    if (old_obj.st == 'starting' && new_st == 'off' && (new_ts - old_obj.ts) < 5 * 60) {
                        new_st = 'starting';
                        new_ts = old_obj.ts;
                    } else if (old_obj.st == 'stopping' && new_st == 'on' && (new_ts - old_obj.ts) < 5 * 60) {
                        new_st = 'stopping';
                        new_ts = old_obj.ts;
                    }
                }

                new_lst[e.id] = {st: new_st, ts: new_ts};
            }

            status_onoff[aid] = new_lst;
        } catch {}
    }
}

setInterval(() => {
    fetch_status_onoff().catch(console.error);
}, 10000);

// ============================================================================

var status_open = {
    // aid: {
    //     opening: 0,
    //     last: {
    //         id: 0,
    //         ts: '',
    //     },
    // },
};

function get_opening(aid) {
    let obj = status_open[aid];
    return obj && obj.opening ? obj.opening : 0;
}

function set_opening(aid, v) {
    let obj = status_open[aid];
    if (!obj) {
        obj = {};
        status_open[aid] = obj;
    }

    obj.opening = v;
}

async function get_last(aid) {
    let obj = status_open[aid];
    if (obj && obj.last) {
        return obj.last;
    } else {
        return await fetch_last(aid);
    }
}

async function fetch_last(aid) {
    let a = area.find(aid);
    if (!a) return {};

    try {
        let doc = await dbpool.exec(a.conf.common.db_center, async db => {
            return await db.collection('lastopen').findOne({_id: 1});
        });

        let obj = status_open[aid];
        if (!obj) {
            obj = {};
            status_open[aid] = obj;
        }
        if (!obj.last) obj.last = {};

        if (doc) {
            obj.last.id = doc.id;
            obj.last.ts = doc.ts;
        }

        return obj.last;
    } catch {
        return {};
    }
}

// timer: fix last in case of script operations outside this program
setInterval(() => {
    (async () => {
        for (let aid of area.ids()) {
            let opening = get_opening(aid);
            if (!opening) {
                await fetch_last(aid);
            }
        }
    })().catch(console.error);
}, 5 * 60 * 1000);

// ============================================================================

module.exports = {
    onoff: status_onoff,

    get_opening: get_opening,
    set_opening: set_opening,
    get_last:    get_last,
    fetch_last:  fetch_last,
}
