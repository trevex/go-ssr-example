// Unfortunately `history.pushState` and `history.replaceState` do not emit
// events, so we wrap them to be able to add classes based on current path
// in our menu.
var wrapHistoryFunc = function(type) {
    var orig = history[type];
    return function() {
        var rv = orig.apply(this, arguments);
        var e = new Event(type.toLowerCase());
        e.arguments = arguments;
        window.dispatchEvent(e);
        return rv;
    };
};
history.pushState = wrapHistoryFunc('pushState');
history.replaceState = wrapHistoryFunc('replaceState');

import htmx from 'htmx.org';
window.htmx = htmx;

import Alpine from 'alpinejs'
window.Alpine = Alpine
Alpine.start()


