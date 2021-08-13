
var express = require('express');
var router  = express.Router();

// ============================================================================

router.use('/summary',      require('./summary'));
router.use('/ltv',          require('./ltv'));
router.use('/rank',         require('./rank'));
router.use('/oldnew',       require('./oldnew'));

// ============================================================================

module.exports = router;
