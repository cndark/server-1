
var express = require('express');
var router  = express.Router();

// ============================================================================

router.use('/summary',     require('./summary'));
router.use('/live',        require('./live'));
router.use('/inst',        require('./inst'));
router.use('/model',       require('./model'));
router.use('/cstep',       require('./cstep'));
router.use('/feedback',    require('./feedback'));



// ============================================================================

module.exports = router;
