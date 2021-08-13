
var express = require('express');
var router  = express.Router();

var fs      = require('fs');
var path    = require('path');
var md5     = require('md5');

// ============================================================================

var CCS_DIR        = null;
var CCS_SIGN_KEY   = "gKh893b7cR2bnfoIgvqc6CgqtdvS67f2";

// ============================================================================

(function () {
    CCS_DIR = path.join(__dirname, "../../../ccs");

    if (!fs.existsSync(CCS_DIR)) {
        try { fs.mkdirSync(CCS_DIR, 0755); } catch(e) {}
    }
})();

// ============================================================================

router.post('/', function (req, res) {
    let seed = req.body.seed;
    let txt  = req.body.text;
    let sign = req.body.sign;

    if (!seed || !txt || !sign) {
        res.status(404).end();
        return;
    }

    // check signature
    let sign2 = md5(`${seed}${txt.replace(/[\r\n]/g, "")}${CCS_SIGN_KEY}`);
    if (sign != sign2) {
        res.json({err: 1, msg: "sign error"}).end();
        return;
    }

    // accepted
    res.json({err: 0}).end();

    // ok. save the callstack
    let fn = path.join(CCS_DIR, md5(txt));

    fs.open(fn, "wx", 0644, (err, fd) => {
        if (err) {
            fs.utimes(fn, new Date(), new Date(), (err1)=>{
                if (err1) return;
            });
        }else{
            fs.write(fd, txt, 0, 'utf8', (err) => {
                fs.close(fd, err2=>console.error('close file failed:', err2));
                if (err) console.error("save ccs failed:", err);
            });
        }
    });
});

// ============================================================================

module.exports = router;
