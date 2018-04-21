var final_transcript = '';
var recognizing = false;
var ignore_onend;
var start_timestamp;

if (!('webkitSpeechRecognition' in window)) {
    upgrade();
} else {
    var recognition = new webkitSpeechRecognition();
    recognition.continuous = true;
    recognition.interimResults = true;

    recognition.onstart = function() {
        recognizing = true;
    };

    recognition.onerror = function(event) {
        if (event.error == 'no-speech') {
            ignore_onend = true;
        }
        if (event.error == 'audio-capture') {
            ignore_onend = true;
        }
        if (event.error == 'not-allowed') {
            ignore_onend = true;
        }
    };

    recognition.onend = function() {
        recognizing = false;
        if (ignore_onend) {
            return;
        }
        if (!final_transcript) {
            return;
        }
        if (window.getSelection) {
            window.getSelection().removeAllRanges();
            var range = document.createRange();
            range.selectNode(document.getElementById('final_span'));
            window.getSelection().addRange(range);
        }
    };

    recognition.onresult = function(event) {
        var interim_transcript = '';
        if (typeof(event.results) == 'undefined') {
            recognition.onend = null;
            recognition.stop();
            upgrade();
            return;
        }
        for (var i = event.resultIndex; i < event.results.length; ++i) {
            if (event.results[i].isFinal) {
                final_transcript += event.results[i][0].transcript;
            } else {
                interim_transcript += event.results[i][0].transcript;
                document.getElementById("results-viewer").textContent=interim_transcript;
            }
        }
        document.getElementById("results-viewer").textContent=final_transcript;
        final_span.innerHTML = final_transcript;
        interim_span.innerHTML = interim_transcript;
    };
}

function startButton(event) {
    if (recognizing) {
        recognition.stop();
        document.getElementById("voice-button").textContent="Processing your command...";
        document.getElementById("voice-button").classList.remove("recording");
        document.getElementById("transcript").setAttribute('value', document.getElementById("final_span").textContent);
        document.getElementById("audioForm").submit();
        return;
    }
    final_transcript = '';
    recognition.lang = 'en-US';
    recognition.start();
    ignore_onend = false;
    final_span.innerHTML = '';
    interim_span.innerHTML = '';
    start_timestamp = event.timeStamp;
    document.getElementById("voice-button").textContent="When finished, click me!";
    document.getElementById("voice-button").classList.add("recording");
    document.getElementById("results-viewer").textContent="Listening...";
}

window.addEventListener('keypress', function (e) {
    if (e.keyCode == 32) {
        startButton(e);
    }
}, false);