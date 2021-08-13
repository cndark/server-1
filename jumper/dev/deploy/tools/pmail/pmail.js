#!/usr/bin/env node

require('../lib/ext');

var axios = require('axios');

var dbpool = require('../lib/dbpool');
var config = require('./config.json');

var args = process.argv.slice(2);
if (args.length < 1) {
    console.log("args: pmail-json-file");
    process.exit(1);
}

/*

pmail-json-file format:

{
    title: '',
    text:  '',

    mails: [
        {
            uid:  '',
            res_k: [],
            res_v: [],
        },

        ...
    ]
}

*/

var pmail = require(`./${args[0]}`);

(async () => {
    try {
        // wait
        console.log('pmail sending will BEGIN in 15 seconds ...');
        await new Promise(rv=>setTimeout(rv, 15000));

        // lookup player svr
        console.log('filling user svr ...');

        await dbpool.exec(config.common.db_share, async db => {
            await Promise.limit(5, pmail.mails, async m => {
                let doc = await db.collection('userinfo').findOne({_id: m.uid});
                if (doc) m.svr = doc.svr;
            });
        });

        // send mails
        console.log('sending mails ...');

        for (let m of pmail.mails) {
            let gs = config.games[m.svr];
            if (!gs) throw `${m.svr} NOT found for ${m.uid}`;

            try {
                await axios.post(`http://${gs.svc}/gm/service`, {
                    plrid: m.uid,
                    key:   "w.pmail",
                    title: pmail.title,
                    text:  pmail.text,
                    res_k: m.res_k,
                    res_v: m.res_v,
                });

                console.log(`${m.uid} OK`);

            } catch(e) {
                console.error(`Error: ${m.uid} -> `, e.message);
            }
        }
    } catch(e) {
        console.error('ERROR:', e);
    }
})().catch(console.error);
