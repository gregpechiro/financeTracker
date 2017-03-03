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
    return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});
}
