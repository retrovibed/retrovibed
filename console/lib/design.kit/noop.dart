void fnNoop<T>(T v) {}

List<T> fnOnChange<T>(List<T> s, T? v, bool Function(T a) cmp) {
  if (v == null) {
    return s.where((o) => !cmp(o)).toList();
  }

  return s.map((o) => cmp(o) ? v : o).toList();
}
