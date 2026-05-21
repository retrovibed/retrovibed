import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/httpx.dart' as httpx;
export 'authn/authenticated.dart';
export 'authn/login.dart';
export 'authn/cache.dart';
export 'authn/deeppool.cache.dart';
export 'authn/api.dart';

Future<String> bearer<T>(authz.Cached<T> c) {
  return c.auto().then((v) => v.bearer);
}

httpx.Option request<T>(authz.Cached<T> c) {
  return httpx.Request.bearer(() => bearer(c));
}
