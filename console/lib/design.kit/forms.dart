export './forms/field.dart';
export './forms/container.dart';
export './forms/input.checkbox.dart';
export './forms/item.management.dart';
export './forms/presets.dart';

typedef FnOnChange<T> = Future<T> Function(T v);

Future<T> FnOnChangeNoop<T>(T v) => Future.value(v);
