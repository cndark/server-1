
$(function ()
{
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
	}

	// fetch init data: res
	{
		$.ajax({
			method: "post",
			url:    "/gift/gtab",
			data:   "key=res",
			dataType: "json",
		}).done(function (r) {
			let combo = $('.panel-body input[name=res_k]');

			combo.typeahead({
                source: r.map(e => `${e[1]} - ${e[0]}`),
            });

			combo.on('change', function () {
				let $this = $(this);
                let e = $this.typeahead('getActive');
                if (!e || e != $this.val())
                    $this.val("");
			});
		});
	}

	// button: gen
	$('#btn_gen').on('click', function () {
		let $this = $(this);
		let ctx = $this.parent().next();

		$.ajax({
			method: "post",
			url: "/gift/gen",
			data: ctx.find('select, input').serialize(),
            dataType: "json",
		}).done(function (r) {
			if (err_show(r)) {
				ctx.find('textarea').val(r.codes.join('\n'));
			}
		}).fail(function () {
			err_show(ErrNet);
		});
	});

	// button: load
	$('#btn_load').on('click', function () {
		let $this = $(this);
		let ctx = $this.parent().next();

		$.ajax({
			method: "post",
			url: "/gift/gen/load",
			data: ctx.find('input[name=grpid]').serialize(),
            dataType: "json",
		}).done(function (r) {
			if (err_show(r)) {
				ctx.find('[name=count]').val(r.codes.length);
				ctx.find('[name=area]').val(r.area);
				ctx.find('[name=reuse]').val(r.reuse);
				ctx.find('[name=expire]').val(r.expire);
				ctx.find('[name=memo]').val(r.memo);

				let arr_k = ctx.find('[name=res_k]');
				let arr_v = ctx.find('[name=res_v]');
				
				for (let i = 0; i < 10; i++) {
					let rwd = r.rewards[i];

					$(arr_k.get(i)).val(rwd ? rwd.res_k : '');
					$(arr_v.get(i)).val(rwd ? rwd.res_v : '');
				}

				ctx.find('textarea').val(r.codes.join('\n'));
			}
		}).fail(function () {
			err_show(ErrNet);
		});
	});

});
