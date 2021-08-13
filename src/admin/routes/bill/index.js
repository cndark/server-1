
var express = require('express');
var router  = express.Router();

// ============================================================================

router.use('/order',    require('./order'));

// ============================================================================

module.exports = router;
