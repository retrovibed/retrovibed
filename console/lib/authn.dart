import 'package:flutter/material.dart';
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn/developer.mode.dart';
import 'package:retrovibed/authn/login.dart';

export 'authn/authenticated.dart';
export 'authn/login.dart';
export 'authn/cache.dart';
export 'authn/deeppool.cache.dart';
export 'authn/api.dart';
export 'authn/developer.mode.dart';

Future<String> bearer<T>(authz.Cached<T> c) {
  return c.auto().then((v) => v.bearer);
}

httpx.Option request<T>(authz.Cached<T> c) {
  return httpx.Request.bearer(() => bearer(c));
}

DeveloperMode developer(BuildContext context) {
  return Login.cached(context).flags;
}
