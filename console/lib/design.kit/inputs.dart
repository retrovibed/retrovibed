export 'inputs/bytes.dart';
export 'inputs/date.dart';
export 'inputs/date.range.dart';
export 'inputs/duration.dart';
export 'inputs/rate.limit.dart';
export 'inputs/time.range.dart';

void _notimplemented(String s) {
  return print(s);
}

void defaulttap() => _notimplemented("tap not implemented");

Future<T> defaulttapfn<T>(T v) {
  _notimplemented("tap not implemented");
  return Future.value(v);
}

Future<void> defaulttapv({String msg = "tap not implemented"}) => Future.sync(() => _notimplemented(msg));
