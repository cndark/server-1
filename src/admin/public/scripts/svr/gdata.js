
function fill_select(name, r) {
	let opt = [];
	for (let row of r) {
		opt.push('<option value="' + row[0] + '">' + row[1] + '</option>');
	}
	$(`.panel-body select[name=${name}]`).html(opt.join());
}

// ============================================================================

$(function ()
{
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
	}

	// fetch init data: conf
	{
		$.ajax({
			method: "post",
			url:    "/server/gtab",
			data:   "key=conf",
			dataType: "json",
		}).done(function (r) {
			fill_select("conf_k", r);
		});
	}

	// normal button event
	$('.head button').on('click', function () {
		let $this = $(this);

		$.ajax({
			method: "post",
			url: '/server/gdata',
			data: $this.closest('.panel').find('select, input, textarea').serialize(),
            dataType: "json",
		}).done(function (r) {
			err_show(r, $this);
		}).fail(function () {
			err_show(ErrNet, $this);
		});
	});

});
