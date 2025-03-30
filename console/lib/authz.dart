import 'package:synchronized/synchronized.dart' as sync;

mixin MyMixin {
  void myMixinMethod() {
    print('Mixin method called');
  }
}

class Cached<T> {
  sync.Lock _m = sync.Lock();

  T current;
  Future<T> Function(Cached<T> current) refresh;

  Cached(this.current, this.refresh);

  Future<T> token() {
    return refresh(this).then(
      (v) => _m.synchronized(() {
        this.current = v;
        return v;
      }),
    );
  }
}

Future<T> Function(Cached<T>) refresh<T>(
  Future<T> Function(T current) fn,
  bool Function(T current, DateTime ts) expired,
) {
  return (t) {
    final ts = DateTime.now();

    if (!expired(t.current, ts)) {
      return Future.value(t.current);
    }

    return fn(t.current);
  };
}
