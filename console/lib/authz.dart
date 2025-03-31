import 'package:synchronized/synchronized.dart' as sync;

class Bearer<T> {
  T metadata;
  String bearer;
  Bearer(this.metadata, this.bearer);
}

class Cached<T> {
  sync.Lock _m = sync.Lock();

  Bearer<T> current;
  Future<Bearer<T>> Function(Cached<T> current) refresh;

  Cached(this.current, this.refresh);

  Future<Bearer<T>> token() {
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

    if (!expired(t.current.metadata, ts)) {
      return Future.value(t.current);
    }

    return fn(t.current.metadata);
  };
}
