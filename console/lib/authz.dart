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

  // returns a refreshed (if necessary) bearer token. Serializes overlapping
  // calls so a slower, superseded fetch can never resolve after (and clobber
  // the result of) a fetch that started later: the whole fetch-and-write is
  // one critical section, not just the write, so a queued call re-checks
  // expiry against the value the prior call just installed instead of
  // racing it.
  Future<Bearer<T>> auto() {
    return _m.synchronized(() {
      return refresh(this).then((v) {
        this.current = v;
        return v;
      });
    });
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
