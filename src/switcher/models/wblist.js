
var ipaddr = require('ipaddr.js');

var dbpool = require('../lib/dbpool');

// ============================================================================

var arr;
var dict;

// ============================================================================

async function load() {
    arr = {
        w: {
            ips: [],
            devices: [],
        },
        b: {
            ips: [],
            devices: [],
        },
    };
    dict = {
        ips: [], // [[CIDR, tp], ...]
        devices: {}, // {deviceid: tp, ...}
    };

    try {
        let db = dbpool.get('center');

        let doc = await db.collection('wblist').findOne({ _id: 1 });
        if (doc) {
            set('w', doc.w_ips, doc.w_devices, true);
            set('b', doc.b_ips, doc.b_devices, true);
        }

        build_dict();
    } catch (e) {
        console.error('loading wblist failed:', e);
    }
}

function set(tp, ips, devices) {
    if (ips) {
        arr[tp].ips = [];
        ips.forEach(v => {
            v = v.trim();
            if (v == "") return;

            arr[tp].ips.push(v);
        });
    }

    if (devices) {
        arr[tp].devices = [];
        devices.forEach(v => {
            v = v.trim();
            if (v == "") return;

            arr[tp].devices.push(v);
        });
    }
}

function build_dict() {
    // black overrides white

    ['b', 'w'].forEach(tp => {
        arr[tp].ips.forEach(v => {
            try {
                var e = [ipaddr.parseCIDR(v), tp];
                dict.ips.push(e);
            } catch (e) { }
        });
    });

    ['w', 'b'].forEach(tp => {
        arr[tp].devices.forEach(v => {
            dict.devices[v] = tp;
        });
    });
}

function type_ip(ip) {
    if (!ip) return '';

    try {
        // parse ip
        var addr = ipaddr.process(ip);

        // match
        for (let v of dict.ips) {
            if (addr.match(v[0])) return v[1];
        }
    } catch (e) { }

    return '';
}

function type_device(deviceid) {
    if (!deviceid) return '';

    var tp = dict.devices[deviceid];
    return tp ? tp : '';
}

// ============================================================================

module.exports = {
    load: load,
    type_ip: type_ip,
    type_device: type_device,
};
