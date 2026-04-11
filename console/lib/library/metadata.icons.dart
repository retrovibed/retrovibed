import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

Icon archived(String uid, {ds.Defaults? defaults, double size = 24.0}) {
  return toggled<Icon>(
    uid,
    Icon(Icons.upload, size: size),
    Icon(Icons.pending_actions_outlined, size: size, color: defaults?.opaque),
    Icon(Icons.archive_outlined, size: size),
  );
}

Icon archived_trash(String uid, {ds.Defaults? defaults, double size = 24.0}) {
  return toggled<Icon>(
    uid,
    Icon(Icons.upload, size: size),
    Icon(Icons.pending_actions_outlined, size: size, color: defaults?.opaque),
    Icon(Icons.delete_forever_rounded, size: size),
  );
}

Widget sharing(String uid, {ds.Defaults? defaults, double  size = 24.0, Color? color}) {
  final minmax = Icon(Icons.share_outlined, size: size, color: defaults?.opaque);
  final v = Icon(Icons.archive_outlined, size: size);
  return toggled(uid, minmax, minmax, v);
}

T toggled<T extends Widget>(String uid, T zero, T max, T o) {
  switch (uid) {
    case "":
    case "00000000-0000-0000-0000-000000000000":
      return zero;
    case "ffffffff-ffff-ffff-ffff-ffffffffffff":
      return max;
    default:
      return o;
  }
}