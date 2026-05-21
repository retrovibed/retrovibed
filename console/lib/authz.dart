import 'package:synchronized/synchronized.dart' as sync;

class Bearer<T> {
  T token;
  String bearer;
  Bearer(this.token, this.bearer);
}

class Cached<T> {
  static Future<Bearer<T>> pending<T>(Cached<T> old) {
    return Future.delayed(Duration(days: 365));
  }

  static Future<Bearer<T>> noprefresh<T>(Cached<T> old) {
    return Future.value(old.current);
  }

  sync.Lock _m = sync.Lock();

  Bearer<T> current;
  Future<Bearer<T>> Function(Cached<T> current) refresh;

  Cached(this.current, this.refresh);

  // returns a refreshed (if necessary) bearer token.
  Future<Bearer<T>> auto() {
    return refresh(this).then(
      (v) => _m.synchronized(() {
        this.current = v;
        return v;
      }),
    );
  }
}

Future<Bearer<T>> Function(Cached<T>) refresh<T>(
  Future<Bearer<T>> Function(T current) fn,
  bool Function(T current, DateTime ts) expired,
) {
  return (t) {
    final ts = DateTime.now();

    if (!expired(t.current.token, ts)) {
      return Future.value(t.current);
    }
    return fn(t.current.token);
  };
}
