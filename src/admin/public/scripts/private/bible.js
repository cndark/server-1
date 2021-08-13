
$(function ()
{
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
    }

    // button: dyncode
	$('#btn_dyncode').on('click', function () {
		let $this = $(this);

		$.ajax({
			method: "post",
			url: '/private/bible/dyncode/set',
            data: $this.closest('.panel').find('select, input, textarea').serialize(),
            dataType: "json",
		}).done(function (r) {
			err_show(r, $this);
		}).fail(function () {
			err_show(ErrNet, $this);
		});
	});
	
	// button: dyncode_gen
	$('#btn_dyncode_gen').on('click', function () {
		let dict = 'abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ23456789';
		let L = dict.length;

		let r = [];
		for (let i = 0; i < 50; i++) {
			r[i] = dict[Math.floor(Math.random() * L)];
		}
		let code = r.join('');

		$(this).closest('.panel').find('input[name=code]').val(code);
	});
});
