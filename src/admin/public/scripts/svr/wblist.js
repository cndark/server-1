
function load() {
	$.ajax({
		method:   "get",
		url:      "/server/wblist/load",
		data:     {area: $('select[name=area]').val()},
		dataType: "json",
	}).done(function (r) {
		if (!err_ok(r)) {
			err_show(ErrNet);
			return;
		}

		let ctx = $('.panel');

		ctx.find('textarea[name=w_ips]').val(r.w_ips.join("\n"));
		ctx.find('textarea[name=w_devices]').val(r.w_devices.join("\n"));

		ctx.find('textarea[name=b_ips]').val(r.b_ips.join("\n"));
		ctx.find('textarea[name=b_devices]').val(r.b_devices.join("\n"));
	});
}

function save() {
	$.ajax({
		method: "post",
		url: '/server/wblist/save',
		data: $('.panel').find('select, textarea').serialize(),
		dataType: "json",
	}).done(function (r) {
		err_show(r);
	}).fail(function () {
		err_show(ErrNet);
	});
}

// ============================================================================

$(function ()
{
	// add .form-control
	{
		$('.panel-body select, input, textarea').addClass("form-control");
	}

	// init load
	{
		load();
	}

	// area change
	$('select[name=area]').on('change', function () {
		load();
	});

	// button: wblist
	$('.head button').on('click', function () {
		save();
	});

});
