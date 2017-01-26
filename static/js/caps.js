var caps = false;
$('#password').keypress(function(e) {
    var s = String.fromCharCode( e.which );
    if ( (s.toUpperCase() === s && s.toLowerCase() !== s && !e.shiftKey) || (s.toLowerCase() === s && s.toUpperCase() !== s && e.shiftKey) ) {
        $('#capsWarn').removeClass('hide');
        caps = true;
    }
});

function onKeyDown(e) {
    if (e.keyCode == 20 && caps) {
        caps = false;
        $('#capsWarn').addClass('hide');
    }
}

document.addEventListener('keydown', onKeyDown, false);
