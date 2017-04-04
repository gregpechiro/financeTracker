function roundTo(num, place) {
	place = Math.pow(10, place);
	return (Math.round(num * place)) / place;
}

$(document).ready(function() {
    $('span.calander').click(function(){
        $(this).parent().find('input').focus()
    });
});

function toTitleCase(str) {
	if (str == '' || str == null) {
		return '';
	}
    return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});
}

var currencyFormatter = new Intl.NumberFormat('en-US', {
 	style: 'currency',
 	currency: 'USD',
 	minimumFractionDigits: 2,
});
