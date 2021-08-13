#!/usr/bin/env node

require('../lib/ext');

var axios = require('axios');

var dbpool = require('../lib/dbpool');
var config = require('./config.json');

// ============================================================================

function _info() {
    console.log.call(this, "[Info]", ...arguments);
}

function _notice() {
    console.log.call(this, "\033[32m[Notice]", ...arguments, "\033[0m");
}

function _warning() {
    console.log.call(this, "\033[33m[Warning]", ...arguments, "\033[0m");
}

function _error() {
    console.log.call(this, "\033[31m[Error]", ...arguments, "\033[0m");
}

function usage() {
    _notice("args: [--check] names...");
    _notice("    this will merge name1, name2, ... into name0");
    process.exit(1);
}

// ============================================================================

var just_check = false;
var games      = [];

process.argv.slice(2).forEach(v => {
    if (v == "--check") {
        just_check = true;
        return;
    }

    var gs = config.games[v];
    if (!gs) {
        _error("game NOT found:", v);
        usage();
    }

    games.push({name: v, db: gs.db_game});
});

if (games.length < 2) {
    usage();
}

if (just_check)
    process.exit(0);

var gs0 = games[0].name;
var gsX = games.slice(1).map(v => v.name);

// ============================================================================

var dbu = Object.keys(config.common.db_user).map(v=>config.common.db_user[v]);
var dbs = config.common.db_share;
var dbc = config.common.db_center;

// ============================================================================

(async () => {
    try {
        // wait  ----------------------------------------
        _notice(`merge ${gsX.join(",")} to ${gs0} will BEGIN in 15 seconds ... press Ctrl + C to quit`);
        await new Promise(rv=>setTimeout(rv, 15000));

        // dbc: update player merge history ----------------------------------------
        _info(`updating merge history ...`);

        {
            // first, get userinfo docs
            let docs = await dbpool.exec(dbs, async db => {
                return await db.collection('userinfo').find({svr: {$in: gsX}}).project({svr: 1}).toArray();
            });

            // then, merge history
            await dbpool.exec(dbc, async db => {
                var bulk = db.collection('mergehis').initializeOrderedBulkOp();

                docs.forEach(doc => {
                    bulk.find({_id: doc._id}).upsert().updateOne({$push: {his: doc.svr}});
                });

                if (bulk.length == 0) {
                    throw "merge data NOT found";
                }

                await bulk.execute();
            });
        }

        // dbgs: remove some data in gs0 ----------------------------------------
        _info(`updating ${gs0} db ...`);

        let gmailid;

        await dbpool.exec(games[0].db, async db => {
            // table removal
            for (let c of ['arenarank']) {
                try { await db.collection(c).drop() } catch {}
            }

            // worlddata record removal
            await db.collection('worlddata').deleteMany(
                {_id: {$in: [
                    "arenaaward",
                ]}},
            );

            // get current gmailid
            let doc = await db.collection('gmail').findOne({_id: 1}, {seqid: 1});
            gmailid = doc ? doc.seqid : 0;
        });

        // dbu: user ----------------------------------------
        _info("updating users ...");

        await Promise.limit(10, dbu, async cnnstr => {
            await dbpool.exec(cnnstr, async db => {
                // update svr & gmailid
                await db.collection('user').updateMany(
                    {"base.svr": {$in: gsX}},
                    {$set: {
                        "base.svr":        gs0,
                        "mailbox.gmailid": gmailid,
                    }},
                );

                // update some data (!Note!: gs0 now includes gsX)
                await db.collection('user').updateMany(
                    {"base.svr": gs0},
                    {$unset: {
                        "base.arena": 1,
                    }},
                );
            });
        });

        // dbu: guild ----------------------------------------
        _info("updating guilds ...");

        await Promise.limit(10, dbu, async cnnstr => {
            await dbpool.exec(cnnstr, async db => {
                await db.collection('guild').updateMany(
                    {"base.svr": {$in: gsX}},
                    {$set: {"base.svr": gs0}},
                );
            });
        });

        // dbs: update userinfo ----------------------------------------
        _info("updating userinfo ...");

        await dbpool.exec(dbs, async db => {
            // userinfo
            await db.collection('userinfo').updateMany(
                {svr: {$in: gsX}},
                {$set: {svr: gs0}},
            );
        });

        // dbc: update svrlist ----------------------------------------
        _info("updating svrlist ...");

        await dbpool.exec(dbc, async db => {
            // svrlist
            await db.collection('svrlist').updateMany(
                {mname: {$in: gsX}},
                {$set: {mname: gs0}},
            );

            // notify switcher
            try {
                await axios.get(
                    `http://${config.switcher.ip}:${config.switcher.port}/server/list/update?token=${config.switcher.token}`,
                );
            } catch {
                _warning("notify switcher failed. it MUST be done manually!");
            }
        });

        // dbgs: update merge info in gs0 ----------------------------------------
        _info(`updating merge info in ${gs0} db ...`);

        await dbpool.exec(games[0].db, async db => {
            // update merge info
            await db.collection('worlddata').updateOne(
                {_id: "svrts"},
                {
                    $set: {merge_ts:  new Date()},
                    $inc: {merge_cnt: 1},
                },
            );
        });

        // dbc: write merge log ----------------------------------------
        _info(`writing merge log ...`);

        await dbpool.exec(dbc, async db => {
            try {
                await db.collection('mergelog').insertOne({from: gsX, to: gs0, ts: new Date()});
            } catch(e) {
                _warning("WARNING: writing merge log failed:", e);
            }
        });

        // ok ----------------------------------------
        _notice("Successfully Done");

    } catch(e) {
        _error("ERROR:", e);
    }

})().catch(console.error);
