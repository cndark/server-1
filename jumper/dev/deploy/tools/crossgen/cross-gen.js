#!/usr/bin/env node

var config = require('./config.json');

// ============================================================================

function usage() {
    console.log("args: N sections");
    console.log("");
    console.log("   N:        N servers each group");
    console.log("   sections: comma-separated numbers for section list.  (e.g. 50,100,150,200,300,500,1000)");
    console.log("");
    process.exit(1);
}

let args = process.argv.slice(2);
if (args.length < 2) {
    usage();
}

let N = Number(args[0]);
let section_str = args[1];

// ============================================================================

let scts = section_str.split(/,+/).map(v => Number(v));

let r    = [];
let p    = [];
let sidx = 0;
let maxid = Object.keys(config.games).map(v => Number(v.match(/game(\d+)/)[1])).sort((a, b) => b - a)[0];
let lastid = scts[scts.length - 1];

for (let id = 1; id <= lastid; id++) {
    if (id <= maxid && !config.games["game" + id]) {
        continue;
    }

    if (id > scts[sidx]) {
        r.push(p);
        p = [];
        sidx++;
        if (sidx >= scts.length) {
            console.error("section list MAY be incorrect!");
            process.exit(1);
        }
    } else if (p.length >= N) {
        r.push(p);
        p = [];
    }

    p.push(id);
}

if (p.length > 0) {
    r.push(p);
    p = [];
}

r = r.map(a => a.join(",")).join("|");
console.log(r);
