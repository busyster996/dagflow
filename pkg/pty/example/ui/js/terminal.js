(function() {
    var terminal = new Terminal({
        screenKeys: true,
        useStyle: true,
        cursorBlink: true,
        fullscreenWin: true,
        maximizeWin: true,
        screenReaderMode: true,
        cols: 128,
        fontSize: 12,
    });
    terminal.open(document.getElementById("terminal"));
    var protocol = (location.protocol === "https:") ? "wss://" : "ws://";
    var url = protocol + location.host + "/ws"
    var ws = new WebSocket(url);
    var attachAddon = new AttachAddon.AttachAddon(ws);
    var fitAddon = new FitAddon.FitAddon();
    terminal.loadAddon(fitAddon);
    var webLinksAddon = new WebLinksAddon.WebLinksAddon();
    terminal.loadAddon(webLinksAddon);
    var unicode11Addon = new Unicode11Addon.Unicode11Addon();
    terminal.loadAddon(unicode11Addon);
    var serializeAddon = new SerializeAddon.SerializeAddon();
    terminal.loadAddon(serializeAddon);
    ws.onclose = function(event) {
        console.log(event);
        terminal.write('\r\n\nconnection has been terminated from the server-side (hit refresh to restart)\n')
    };
    ws.onopen = function() {
        terminal.loadAddon(attachAddon);
        terminal._initialized = true;
        terminal.focus();
        setTimeout(function() {
            fitAddon.fit();
        });
        terminal.onResize(function(event) {
            var rows = event.rows*12;
            var cols = event.cols*4;
            var size = JSON.stringify({cols: cols, rows: rows});
            var send = new TextEncoder().encode("\x01" + size);
            console.log('resizing to', size);
            ws.send(send);
            fitAddon.fit();
        });
        terminal.onTitleChange(function(event) {
            console.log(event);
        });
        window.onresize = function() {
            fitAddon.fit();
        };
    };
})();