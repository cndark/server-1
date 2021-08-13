#!/usr/bin/env node

hbars = require('handlebars');
fs = require('fs');

var args = process.argv.slice(2);
var ctx = {
    svr:   Number(args[0]),
    batch: Number(args[1]),
};

async function gen() {
    hbars.registerHelper('format', function(s, ...p) {
        p.pop();
        let i = 0;
        return s.replace(/%s/g, () => p[i++]);
    });
    hbars.registerHelper('add', function(...a) {
        a.pop();
        if (a.length == 0) return;

        let tp = typeof a[0];
        if (tp == 'string') {
            return a.join('');
        } else if (tp == 'number') {
            return a.reduce((a, v) => a + Number(v));
        }
    });
    hbars.registerHelper('mul', function(...a) {
        a.pop();
        if (a.length == 0) return;

        let tp = typeof a[0];
        if (tp == 'string') {
            return 0;
        } else if (tp == 'number') {
            return a.reduce((a, v) => a * Number(v));
        }
    });


    let tpl = await fs.promises.readFile('./bot.tpl', 'utf8');
    let r = hbars.compile(tpl, {noEscape: true})(ctx);
    console.log(r);
}

gen().catch(console.error);
