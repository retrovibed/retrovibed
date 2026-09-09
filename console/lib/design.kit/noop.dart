void fnNoop<T>(T v) {}

Future<void> fnAsyncNoop() {
  return Future.value();
}

Future<void> fnAsyncNoopOnChangeV2<T>(T orig, T upd) {
  return Future.value();
}

Future<void> fnAsyncNoopOnDelete<T>(T orig) {
  return Future.value();
}

Future<void> fnAsyncNoopOnChange<T>(T v) {
  return Future.value();
}

List<T> fnOnChange<T>(List<T> s, T? v, bool Function(T a) cmp) {
  if (v == null) {
    return s.where((o) => !cmp(o)).toList();
  }

  return s.map((o) => cmp(o) ? v : o).toList();
}
