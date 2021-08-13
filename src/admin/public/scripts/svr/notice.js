
function notice_load() {
	let $btn = $('#btn_notice');
	let ctx = $btn.parent().next();

	$.ajax({
		method: "post",
		url: "/server/notice/load",
		data: ctx.find('select').serialize(),
        dataType: "json",
	}).done(function (r) {
		if (!err_ok(r)) {
			r.d_start = "";
			r.d_end   = "";
			r.title   = "";
			r.content = "";
		}

		ctx.find('[name=d_start]').val(r.d_start);
		ctx.find('[name=d_end]').val(r.d_end);
		ctx.find('[name=title]').val(r.title);
		ctx.find('[name=content]').val(r.content);
	}).fail(function () {
        err_show(ErrNet, $btn);
    });
}

function notice_save() {
	let $btn = $('#btn_notice');
	let ctx = $btn.parent().next();

	$.ajax({
		method: "post",
		url: "/server/notice/save",
		data: ctx.find('select, input, textarea').serialize(),
        dataType: "json",
	}).done(function (r) {
		err_show(r, $btn);
	}).fail(function () {
		err_show(ErrNet, $btn);
	});
}

// ============================================================================

$(function () {
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
	}

	// notice: init load
	{
		notice_load();
	}

	// notice: area/lang change
	$('#btn_notice').parent().next().find('select').on('change', function () {
		notice_load();
	});

	// button: notice
	$('#btn_notice').on('click', function () {
		notice_save();
	});

	// button: lamp
	$('#btn_lamp').on('click', function () {
		let $this = $(this);
		let ctx = $this.parent().next();

		$.ajax({
			method: "post",
			url: "/server/notice/lamp",
			data: ctx.find('select, input, textarea').serialize(),
			dataType: "json",
		}).done(function (r) {
			err_show(r, $this);
		}).fail(function () {
			err_show(ErrNet, $this);
		});
	});

});
