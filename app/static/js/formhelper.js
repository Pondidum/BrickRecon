Turbolinks.savedScrolls = {};

document.addEventListener("turbolinks:before-visit", function (event) {
  Turbolinks.savedScrolls = {
    [window.location.href]: {
      document: document.documentElement.scrollTop,
      body: document.body.scrollTop,
    },
  };
});

document.addEventListener("turbolinks:render", function (event) {
  const savedScroll = Turbolinks.savedScrolls[window.location.href];
  if (!savedScroll) {
    return;
  }

  delete Turbolinks.savedScrolls[window.location.href];

  if (savedScroll.document != null) {
    if (savedScroll.document < document.documentElement.scrollHeight) {
      document.documentElement.scrollTop = savedScroll.document;
    }
  }

  if (savedScroll.document != null) {
    if (savedScroll.body < document.body.scrollHeight) {
      document.body.scrollTop = savedScroll.body;
    }
  }
});

function urlencodeFormData(fd) {
  var s = "";
  function encode(s) {
    return encodeURIComponent(s).replace(/%20/g, "+");
  }
  for (var pair of fd.entries()) {
    if (typeof pair[1] == "string") {
      s += (s ? "&" : "") + encode(pair[0]) + "=" + encode(pair[1]);
    }
  }
  return s;
}

function submitForm(form, e) {
  var xhr = new XMLHttpRequest();
  xhr.onload = () => eval(xhr.responseText);

  xhr.open(
    "POST",
    location.protocol +
      "//" +
      location.host +
      location.pathname +
      location.search
  );

  xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
  xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

  var data = new FormData(form);
  data.append(e.submitter.name, e.submitter.value);
  xhr.send(urlencodeFormData(data));

  return false;
}
