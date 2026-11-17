// CloudFront viewer-request function (runtime cloudfront-js-2.0).
//
// The site is built with trailingSlash: 'always', and a private S3 bucket behind origin
// access control does not serve directory indexes, so:
//   /              -> /index.html
//   /commands/     -> /commands/index.html
//   /commands      -> 301 redirect to /commands/ (query kept)
//   /favicon.svg, /_astro/x.css (last segment has a dot) -> unchanged
//
// Adapted from ninegene/ninegene.com (infra/pulumi/functions/viewer-request.js), without
// its www redirect: this site is a subdomain with no www variant.
// Keep this file free of imports: CloudFront runs it as-is, and the tests in
// viewer-request.test.mjs load it as plain text. CloudFront calls handler by name.

// eslint-disable-next-line no-unused-vars
function handler(event) {
  var request = event.request;
  var uri = request.uri;

  if (uri.charAt(uri.length - 1) === '/') {
    request.uri = uri + 'index.html';
    return request;
  }

  var lastSegment = uri.substring(uri.lastIndexOf('/') + 1);
  if (lastSegment.indexOf('.') === -1) {
    // Collapse leading slashes and backslashes: a Location of //host/ or /\host/
    // is protocol-relative, so browsers would leave the site (open redirect).
    var path = uri.replace(/^[\/\\]+/, '/');
    return redirect(path + '/' + queryString(request.querystring));
  }

  return request;
}

function redirect(location) {
  return {
    statusCode: 301,
    statusDescription: 'Moved Permanently',
    headers: { location: { value: location } },
  };
}

function queryString(query) {
  var parts = [];
  for (var name in query) {
    var entry = query[name];
    if (entry.multiValue) {
      for (var i = 0; i < entry.multiValue.length; i++) {
        parts.push(name + '=' + entry.multiValue[i].value);
      }
    } else {
      parts.push(entry.value === '' ? name : name + '=' + entry.value);
    }
  }
  return parts.length ? '?' + parts.join('&') : '';
}
