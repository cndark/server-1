#!/usr/bin/env node

require('../lib/ext');

var path  = require('path');
var axios = require('axios');

var dbpool = require('../lib/dbpool');

// ============================================================================

function help() {
    console.log(`Usage: ${path.basename(process.argv[1])} [options] id name`);
    console.log();
    console.log("id: range [1, 3000)");
    console.log("")
    console.log("options:");
    console.log("   -s <1|2|3>                  set status");
    console.log("   -f <new|good|test>          set flag");

    process.exit(1);
}

// ============================================================================

let args = process.argv.slice(2);

let opt;
let status = 3;
let flag   = "new";
let rest   = [];

args.forEach(v => {
    if (v == "-s" || v == "-f") {
        opt = v;
    } else if (opt == "-s") {
        status = v;
    } else if (opt == "-f") {
        flag = v;
    } else {
        rest.push(v);
    }
});

if (rest.length < 2) help();

let id   = tonumber(rest[0]);
let name = rest[1];

if (id <= 0 || id >= 3000) help();

// ============================================================================

(async () => {
    // open db center
    let str = await shell_exec(`${path.join(__dirname, './deploy.js')} -c`);
    let c   = JSON.parse(str);

    // insert
    await dbpool.exec(c.common.db_center, async db => {
        let now = new Date();

        try {
            await db.collection('svrlist').insertOne({
                _id:          `game${id}`,
                id:           id,
                mname:        `game${id}`,
                text:         name,
                status:       status,
                flag:         flag,
                closereg:     false,
                tri_hwater:   false,
                tri_closereg: false,
                ts:           now,
            });

            try {
                await db.collection('lastopen').replaceOne(
                    {_id: 1},
                    {id: id, ts: now},
                    {upsert: true},
                );
            } catch(e) {
                console.error('update lastopen failed:', e.message);
            }
        } catch(e) {
            if (e.message.match(/^E11000/)) {
                throw `svrid already exists: ${id}`;
            } else {
                throw e;
            }
        }
    });

    // notify switcher
    try {
        await axios.get(
            `http://${c.switcher.ip}:${c.switcher.port}/server/list/update?token=${c.switcher.token}`,
        );
    } catch(e) {
        console.error('notify switcher failed:', e.message);
    }
})().catch(e => {
    console.error(typeof e == 'string' ? e : e.message);
    process.exit(1);
});
